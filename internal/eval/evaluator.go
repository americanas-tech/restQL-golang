package eval

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/parser"
	"github.com/b2wdigital/restQL-golang/internal/runner"
	"github.com/pkg/errors"
)

var (
	ErrInvalidRevision  = errors.New("revision must be greater than 0")
	ErrInvalidQueryId   = errors.New("query id must be not empty")
	ErrInvalidNamespace = errors.New("namespace must be not empty")
)

type Evaluator struct {
	log            domain.Logger
	mappingsReader MappingsReader
	queryReader    QueryReader
	run            runner.Runner
}

func NewEvaluator(log domain.Logger, mr MappingsReader, qr QueryReader, r runner.Runner) Evaluator {
	return Evaluator{log: log, mappingsReader: mr, queryReader: qr, run: r}
}

func (e Evaluator) SavedQuery(ctx context.Context, queryOpts domain.QueryOptions, queryInput domain.QueryInput) (domain.Resources, error) {
	err := validateQueryOptions(queryOpts)
	if err != nil {
		return nil, err
	}

	queryTxt, err := e.queryReader.Get(nil, queryOpts.Namespace, queryOpts.Id, queryOpts.Revision)
	if err != nil {
		return nil, err
	}

	query, err := parser.Parse(queryTxt)
	if err != nil {
		e.log.Debug("failed to parse query", "error", err)
		return nil, ParserError{errors.Wrap(err, "invalid query syntax")}
	}

	mappings := e.mappingsReader.FromTenant(ctx, queryOpts.Tenant)

	queryCtx := domain.QueryContext{
		Mappings: mappings,
		Options:  queryOpts,
		Input:    queryInput,
	}

	resources, err := e.run.ExecuteQuery(ctx, query, queryCtx)
	switch {
	case err == runner.ErrQueryTimedOut:
		return nil, TimeoutError{Err: err}
	case err != nil:
		return nil, err
	}

	resources, err = ApplyFilters(query, resources)
	if err != nil {
		return nil, err
	}

	resources = ApplyAggregators(query, resources)

	return resources, nil
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

	return nil
}
