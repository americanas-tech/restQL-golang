package eval

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/domain/resolvers"
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

func NewEvaluator(mr MappingsReader, qr QueryReader, r runner.Runner, log domain.Logger) Evaluator {
	return Evaluator{log: log, mappingsReader: mr, queryReader: qr, run: r}
}

func (e Evaluator) SavedQuery(ctx context.Context, queryOpts QueryOptions, queryInput QueryInput) (interface{}, error) {
	err := validateQueryOptions(queryOpts)
	if err != nil {
		return domain.Query{}, err
	}

	queryTxt, err := e.queryReader.GetQuery(queryOpts.Namespace, queryOpts.Id, queryOpts.Revision)
	if err != nil {
		return domain.Query{}, err
	}

	query, err := parser.Parse(queryTxt)
	if err != nil {
		e.log.Debug("failed to parse query", "error", err)
		return domain.Query{}, ParserError{errors.Wrap(err, "invalid query syntax")}
	}

	mappings, err := e.fetchMappings(queryOpts.Tenant, query)
	if err != nil {
		return domain.Query{}, err
	}

	query = resolvers.ResolveVariables(query, queryInput.Params)
	query = resolvers.MultiplexStatements(query)

	r := e.run.ExecuteQuery(ctx, query, mappings)

	return r, nil
}

func (e Evaluator) fetchMappings(tenant string, query domain.Query) (map[string]domain.Mapping, error) {
	mappings := make(map[string]domain.Mapping)

	for _, stmt := range query.Statements {
		url, err := e.mappingsReader.GetMapping(tenant, stmt.Resource)
		if err != nil {
			return nil, err
		}

		mappings[stmt.Resource] = url
	}

	return mappings, nil

}

func validateQueryOptions(queryOpts QueryOptions) error {
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
