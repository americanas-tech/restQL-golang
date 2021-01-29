package eval

import (
	"context"
	"fmt"

	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
	"github.com/b2wdigital/restQL-golang/v4/internal/parser"
	"github.com/b2wdigital/restQL-golang/v4/internal/platform/plugins"
	"github.com/b2wdigital/restQL-golang/v4/internal/runner"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
	"github.com/pkg/errors"
)

// Validation errors returned be Evaluator
var (
	ErrInvalidRevision  = errors.New("revision must be greater than 0")
	ErrInvalidQueryID   = errors.New("query id must be not empty")
	ErrInvalidNamespace = errors.New("namespace must be not empty")
	ErrInvalidTenant    = errors.New("tenant must be not empty")
)

// Evaluator is the interpreter of the restQL language.
// It can execute saved or ad-hoc queries.
type Evaluator struct {
	log            restql.Logger
	parser         parser.Parser
	mappingsReader MappingsReader
	queryReader    QueryReader
	runner         runner.Runner
	lifecycle      plugins.Lifecycle
}

// NewEvaluator constructs an instance of the restQL interpreter.
func NewEvaluator(log restql.Logger, mr MappingsReader, qr QueryReader, r runner.Runner, p parser.Parser, l plugins.Lifecycle) Evaluator {
	return Evaluator{
		log:            log,
		mappingsReader: mr,
		queryReader:    qr,
		runner:         r,
		parser:         p,
		lifecycle:      l,
	}
}

// AdHocQuery executes an ad-hoc send by the client with
// the options and HTTP information.
func (e Evaluator) AdHocQuery(ctx context.Context, queryTxt string, queryOpts restql.QueryOptions, queryInput restql.QueryInput) (domain.Resources, error) {
	if queryOpts.Tenant == "" {
		return nil, ValidationError{ErrInvalidTenant}
	}

	return e.evaluateQuery(ctx, queryTxt, queryOpts, queryInput)
}

// SavedQuery executes a saved query identified by namespace,
// id and revision with the options and HTTP information
// send by the client.
func (e Evaluator) SavedQuery(ctx context.Context, queryOpts restql.QueryOptions, queryInput restql.QueryInput) (domain.Resources, error) {
	err := validateQueryOptions(queryOpts)
	if err != nil {
		return nil, err
	}

	savedQuery, err := e.queryReader.Get(ctx, queryOpts.Namespace, queryOpts.Id, queryOpts.Revision)
	if err != nil {
		return nil, err
	}

	log := restql.GetLogger(ctx)
	log.Debug("Saved query retrieved", "query", savedQuery)

	if savedQuery.Deprecated {
		return nil, domain.ErrQueryRevisionDeprecated{Revision: queryOpts.Revision}
	}

	return e.evaluateQuery(ctx, savedQuery.Text, queryOpts, queryInput)
}

func (e Evaluator) evaluateQuery(ctx context.Context, queryTxt string, queryOpts restql.QueryOptions, queryInput restql.QueryInput) (domain.Resources, error) {
	log := restql.GetLogger(ctx)

	query, err := e.parser.Parse(queryTxt)
	if err != nil {
		log.Debug("failed to parse query", "error", err)
		return nil, ParserError{errors.Wrap(err, "invalid query syntax")}
	}

	mappings, err := e.mappingsReader.FromTenant(ctx, queryOpts.Tenant)
	if err != nil {
		log.Error("failed to fetch mappings", err)
		return nil, err
	}

	err = validateQueryResources(query, mappings)
	if err != nil {
		log.Error("query reference invalid resource", err, "mappings", fmt.Sprintf("%#v", mappings))
		return nil, err
	}

	queryContext := restql.QueryContext{
		Mappings: mappings,
		Options:  queryOpts,
		Input:    queryInput,
	}

	queryCtx := e.lifecycle.BeforeQuery(ctx, queryTxt, queryContext)

	query = ResolveVariables(query, queryContext.Input)

	resources, err := e.runner.ExecuteQuery(queryCtx, query, queryContext)
	switch {
	case err == runner.ErrQueryTimedOut:
		return nil, TimeoutError{Err: err}
	case errors.Is(err, runner.ErrInvalidChainedParameter):
		return nil, ParserError{Err: err}
	case err != nil:
		return nil, err
	}

	resources, err = ApplyFilters(log, query, resources)
	if err != nil {
		log.Error("failed to apply filters", err, "input", fmt.Sprintf("%+#v", queryContext.Input))
		return nil, err
	}

	resources = ApplyAggregators(nil, query, resources)

	e.lifecycle.AfterQuery(queryCtx, queryTxt, resources)

	resources = ApplyHidden(query, resources)

	return resources, nil
}

func validateQueryResources(query domain.Query, mappings map[string]restql.Mapping) error {
	for _, s := range query.Statements {
		_, found := mappings[s.Resource]
		if !found {
			err := errors.Errorf("statement should reference a valid mapped resource. Error was in %s", s.Resource)
			return MappingError{Err: err}
		}
	}

	return nil
}

func validateQueryOptions(queryOpts restql.QueryOptions) error {
	if queryOpts.Revision <= 0 {
		return ValidationError{ErrInvalidRevision}
	}

	if queryOpts.Id == "" {
		return ValidationError{ErrInvalidQueryID}
	}

	if queryOpts.Namespace == "" {
		return ValidationError{ErrInvalidNamespace}
	}

	if queryOpts.Tenant == "" {
		return ValidationError{ErrInvalidTenant}
	}

	return nil
}
