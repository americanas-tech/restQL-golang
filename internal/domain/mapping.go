package domain

import "regexp"

var pathParamRegex, _ = regexp.Compile("/:(\\w+)")
var schemaRegex, _ = regexp.Compile("^(\\w+)://(.+)$")

type Mapping struct {
	ResourceName  string
	Schema        string
	Uri           string
	PathParams    []string
	pathParamsSet map[string]struct{}
}

func NewMapping(resource, url string) Mapping {
	paramsMatches := pathParamRegex.FindAllStringSubmatch(url, -1)

	pathParamsSet := make(map[string]struct{})
	pathParams := make([]string, len(paramsMatches))
	for i, m := range paramsMatches {
		paramName := m[1]
		pathParams[i] = paramName
		pathParamsSet[paramName] = struct{}{}
	}

	schemaMatches := schemaRegex.FindAllStringSubmatch(url, -1)
	schema := schemaMatches[0][1]
	uri := schemaMatches[0][2]

	return Mapping{
		ResourceName:  resource,
		Uri:           uri,
		Schema:        schema,
		PathParams:    pathParams,
		pathParamsSet: pathParamsSet,
	}
}

func (m Mapping) HasParam(name string) bool {
	_, found := m.pathParamsSet[name]
	return found
}
