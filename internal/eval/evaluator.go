package eval

type QueryOptions struct {
	Namespace string
	Id        string
	Revision  int
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

func (e Evaluator) SavedQuery(queryOpts QueryOptions) (interface{}, error) {
	var queryConf queryConfig
	err := e.config.File().Unmarshal(&queryConf)
	if err != nil {
		return nil, err
	}

	e.log.Debug("queries in yaml", "queries", queryConf)

	queriesInNamespace := queryConf.Queries[namespace(queryOpts.Namespace)]
	queriesByRevision := queriesInNamespace[queryOpts.Id]
	queryTxt := queriesByRevision[queryOpts.Revision-1]

	return queryTxt, nil
}
