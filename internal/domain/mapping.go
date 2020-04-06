package domain

import (
	"github.com/pkg/errors"
	"regexp"
)

var pathParamRegex, _ = regexp.Compile("/:(.+)")
var schemaRegex, _ = regexp.Compile("^(\\w+)://(.+)$")

type Mapping struct {
	ResourceName  string
	Schema        string
	Uri           string
	PathParams    []string
	PathParamsSet map[string]struct{}
}

func NewMapping(resource, url string) (Mapping, error) {
	paramsMatches := pathParamRegex.FindAllStringSubmatch(url, -1)

	pathParamsSet := make(map[string]struct{})
	pathParams := make([]string, len(paramsMatches))
	for i, m := range paramsMatches {
		paramName := m[1]
		pathParams[i] = paramName
		pathParamsSet[paramName] = struct{}{}
	}

	schemaMatches := schemaRegex.FindAllStringSubmatch(url, -1)
	if len(schemaMatches) == 0 {
		return Mapping{}, errors.Errorf("failed to create mapping from %s", url)
	}

	if len(schemaMatches[0]) != 3 {
		return Mapping{}, errors.Errorf("failed to create mapping from %s", url)
	}

	schema := schemaMatches[0][1]
	uri := schemaMatches[0][2]

	return Mapping{
		ResourceName:  resource,
		Uri:           uri,
		Schema:        schema,
		PathParams:    pathParams,
		PathParamsSet: pathParamsSet,
	}, nil
}

func (m Mapping) HasParam(name string) bool {
	_, found := m.PathParamsSet[name]
	return found
}
