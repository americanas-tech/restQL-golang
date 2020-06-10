package eval

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/parser"
	"github.com/b2wdigital/restQL-golang/internal/platform/plugins"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"github.com/b2wdigital/restQL-golang/pkg/restql"
	"github.com/pkg/errors"
)

var (
	ErrInvalidRevision  = errors.New("revision must be greater than 0")
	ErrInvalidQueryId   = errors.New("query id must be not empty")
	ErrInvalidNamespace = errors.New("namespace must be not empty")
	ErrInvalidTenant    = errors.New("tenant must be not empty")
)

type Evaluator struct {
	log            restql.Logger
	parser         parser.Parser
	mappingsReader MappingsReader
	queryReader    QueryReader
	runner         runner.Runner
	pluginsManager plugins.Manager
}

func NewEvaluator(log restql.Logger, mr MappingsReader, qr QueryReader, r runner.Runner, p parser.Parser, pm plugins.Manager) Evaluator {
	return Evaluator{
		log:            log,
		mappingsReader: mr,
		queryReader:    qr,
		runner:         r,
		parser:         p,
		pluginsManager: pm,
	}
}

func (e Evaluator) AdHocQuery(ctx context.Context, queryTxt string, queryOpts domain.QueryOptions, queryInput domain.QueryInput) (domain.Resources, error) {
	if queryOpts.Tenant == "" {
		return nil, ValidationError{ErrInvalidTenant}
	}

	return e.evaluateQuery(ctx, queryTxt, queryOpts, queryInput)
}

func (e Evaluator) SavedQuery(ctx context.Context, queryOpts domain.QueryOptions, queryInput domain.QueryInput) (domain.Resources, error) {
	err := validateQueryOptions(queryOpts)
	if err != nil {
		return nil, err
	}

	savedQuery, err := e.queryReader.Get(ctx, queryOpts.Namespace, queryOpts.Id, queryOpts.Revision)
	if err != nil {
		return nil, err
	}

	if savedQuery.Deprecated {
		return nil, domain.ErrQueryRevisionDeprecated{Revision: queryOpts.Revision}
	}

	return e.evaluateQuery(ctx, savedQuery.Text, queryOpts, queryInput)
}

func (e Evaluator) evaluateQuery(ctx context.Context, queryTxt string, queryOpts domain.QueryOptions, queryInput domain.QueryInput) (domain.Resources, error) {
	query, err := e.parser.Parse(queryTxt)
	if err != nil {
		e.log.Debug("failed to parse query", "error", err)
		return nil, ParserError{errors.Wrap(err, "invalid query syntax")}
	}

	mappings, err := e.mappingsReader.FromTenant(ctx, queryOpts.Tenant)
	if err != nil {
		e.log.Error("failed to fetch mappings", err)
		return nil, err
	}

	err = validateQueryResources(query, mappings)
	if err != nil {
		e.log.Error("query reference invalid resource", err)
		return nil, err
	}

	queryContext := domain.QueryContext{
		Mappings: mappings,
		Options:  queryOpts,
		Input:    queryInput,
	}

	queryCtx := e.pluginsManager.RunBeforeQuery(ctx, queryTxt, queryContext)

	resources, err := e.runner.ExecuteQuery(queryCtx, query, queryContext)
	switch {
	case err == runner.ErrQueryTimedOut:
		return nil, TimeoutError{Err: err}
	case errors.Is(err, runner.ErrInvalidChainedParameter):
		return nil, ParserError{Err: err}
	case err != nil:
		return nil, err
	}

	resources, err = ApplyFilters(query, resources)
	if err != nil {
		e.log.Error("failed to apply filters", err)
		return nil, err
	}

	resources = ApplyAggregators(query, resources)

	queryCtx = e.pluginsManager.RunAfterQuery(queryCtx, queryTxt, resources)

	resources = ApplyHidden(query, resources)

	return resources, nil
}

func validateQueryResources(query domain.Query, mappings map[string]domain.Mapping) error {
	for _, s := range query.Statements {
		_, found := mappings[s.Resource]
		if !found {
			err := errors.Errorf("statement should reference a valid mapped resource. Error was in %s", s.Resource)
			return MappingError{Err: err}
		}
	}

	return nil
}

func validateQueryOptions(queryOpts domain.QueryOptions) error {
	if queryOpts.Revision <= 0 {
		return ValidationError{ErrInvalidRevision}
	}

	if queryOpts.Id == "" {
		return ValidationError{ErrInvalidQueryId}
	}

	if queryOpts.Namespace == "" {
		return ValidationError{ErrInvalidNamespace}
	}

	if queryOpts.Tenant == "" {
		return ValidationError{ErrInvalidTenant}
	}

	return nil
}
