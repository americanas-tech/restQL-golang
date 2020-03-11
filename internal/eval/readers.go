package eval

import (
	"github.com/pkg/errors"
	"regexp"
)

type Mapping struct {
	ResourceName string
	Url          string
	PathParams   []string
}

var pathParamRegex, _ = regexp.Compile("/:(\\w+)")

func newMapping(resource, url string) Mapping {
	matches := pathParamRegex.FindAllStringSubmatch(url, -1)

	pathParams := make([]string, len(matches))
	for i, m := range matches {
		pathParams[i] = m[1]
	}

	return Mapping{
		ResourceName: resource,
		Url:          url,
		PathParams:   pathParams,
	}
}

type MappingsReader struct {
	env  EnvSource
	file map[string]string
}

func NewMappingReader(config Configuration, log Logger) MappingsReader {
	mr := MappingsReader{env: config.Env()}

	mappingsConf := struct {
		Mappings map[string]string
	}{}

	err := config.File().Unmarshal(&mappingsConf)
	if err != nil {
		log.Debug("failed to load mappings from config file", "error", err)
	} else {
		mr.file = mappingsConf.Mappings
	}

	return mr
}

func (mr MappingsReader) GetMapping(tenant, resource string) (Mapping, error) {
	switch {
	case mr.env.GetString(resource) != "":
		return newMapping(resource, mr.env.GetString(resource)), nil
	case mr.file[resource] != "":
		return newMapping(resource, mr.file[resource]), nil
	default:
		return Mapping{}, NotFoundError{errors.Errorf("resource `%s` not found on mappings", resource)}
	}
}

type savedQueries map[string][]string

type QueryReader struct {
	file map[string]savedQueries
}

func NewQueryReader(config Configuration, log Logger) QueryReader {
	queryConf := struct {
		Queries map[string]savedQueries
	}{}

	err := config.File().Unmarshal(&queryConf)
	if err != nil {
		log.Debug("failed to load queries from config file", "error", err)
	}

	return QueryReader{file: queryConf.Queries}
}

func (qr QueryReader) GetQuery(namespace, id string, revision int) (string, error) {
	queriesInNamespace, ok := qr.file[namespace]
	if !ok {
		return "", NotFoundError{errors.Errorf("namespace not found: %s", namespace)}
	}

	queriesByRevision, ok := queriesInNamespace[id]
	if !ok {
		return "", NotFoundError{errors.Errorf("query not found: %s", id)}
	}

	if len(queriesByRevision) < revision {
		return "", NotFoundError{errors.Errorf("revision not found: %d", revision)}
	}

	queryTxt := queriesByRevision[revision-1]

	return queryTxt, nil
}
