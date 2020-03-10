package eval

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/parser"
	"github.com/pkg/errors"
)

var (
	ErrInvalidRevision  = errors.New("revision must be greater than 0")
	ErrInvalidQueryId   = errors.New("query id must be not empty")
	ErrInvalidNamespace = errors.New("namespace must be not empty")
)

type Evaluator struct {
	log            Logger
	mappingsReader MappingsReader
	queryReader    QueryReader
}

func NewEvaluator(log Logger, mr MappingsReader, qr QueryReader) Evaluator {
	return Evaluator{log: log, mappingsReader: mr, queryReader: qr}
}

func (e Evaluator) SavedQuery(queryOpts QueryOptions, queryInput QueryInput) (domain.Query, error) {
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

	_, err = e.fetchMappings(query)
	if err != nil {
		return domain.Query{}, err
	}

	query = ResolveVariables(query, queryInput)

	return query, nil
}

func (e Evaluator) fetchMappings(query domain.Query) (map[string]string, error) {
	mappings := make(map[string]string)

	for _, stmt := range query.Statements {
		url, err := e.mappingsReader.GetUrl(stmt.Resource)
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
