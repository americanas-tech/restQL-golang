package eval

import (
	"errors"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/parser"
)

var (
	ErrInvalidRevision  = errors.New("revision must be greater than 0")
	ErrInvalidQueryId   = errors.New("query id must be not empty")
	ErrInvalidNamespace = errors.New("namespace must be not empty")
)

type QueryOptions struct {
	Namespace string
	Id        string
	Revision  int
}

type QueryInput struct {
	Params  map[string]interface{}
	Headers map[string]string
}

type namespace string
type savedQueries map[string][]string

type queryConfig struct {
	Queries map[namespace]savedQueries
}

type Evaluator struct {
	config Configuration
	log    Logger
}

func NewEvaluator(config Configuration, log Logger) Evaluator {
	return Evaluator{config: config, log: log}
}

func (e Evaluator) SavedQuery(queryOpts QueryOptions, queryInput QueryInput) (domain.Query, error) {
	err := validateQueryOptions(queryOpts)
	if err != nil {
		return domain.Query{}, err
	}

	var queryConf queryConfig
	err = e.config.File().Unmarshal(&queryConf)
	if err != nil {
		e.log.Debug("failed to load queries from config file", "error", err)
		return domain.Query{}, err
	}

	queriesInNamespace := queryConf.Queries[namespace(queryOpts.Namespace)]
	queriesByRevision := queriesInNamespace[queryOpts.Id]
	queryTxt := queriesByRevision[queryOpts.Revision-1]

	query, err := parser.Parse(queryTxt)
	if err != nil {
		e.log.Debug("failed to parse query", "error", err)
		return domain.Query{}, err
	}

	query = ResolveVariables(query, queryInput)

	return query, nil
}

func validateQueryOptions(queryOpts QueryOptions) error {
	if queryOpts.Revision <= 0 {
		return ErrInvalidRevision
	}

	if queryOpts.Id == "" {
		return ErrInvalidQueryId
	}

	if queryOpts.Namespace == "" {
		return ErrInvalidNamespace
	}

	return nil
}
