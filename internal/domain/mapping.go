package domain

import (
	"github.com/pkg/errors"
	"regexp"
)

var pathParamRegex = regexp.MustCompile("/:(.+)")
var urlRegex = regexp.MustCompile("(https?)://([-a-zA-Z0-9@:%._+~#=]+)([-a-zA-Z0-9@:%_+.~#?&/=]*)")

type Mapping struct {
	ResourceName  string
	Schema        string
	Host          string
	Path          string
	PathParams    []string
	PathParamsSet map[string]struct{}
}

func NewMapping(resource, url string) (Mapping, error) {
	mapping := Mapping{ResourceName: resource}

	paramsMatches := pathParamRegex.FindAllStringSubmatch(url, -1)

	pathParamsSet := make(map[string]struct{})
	pathParams := make([]string, len(paramsMatches))
	for i, m := range paramsMatches {
		paramName := m[1]
		pathParams[i] = paramName
		pathParamsSet[paramName] = struct{}{}
	}

	mapping.PathParams = pathParams
	mapping.PathParamsSet = pathParamsSet

	urlMatches := urlRegex.FindAllStringSubmatch(url, -1)
	if len(urlMatches) == 0 {
		return Mapping{}, errors.Errorf("failed to create mapping from %s", url)
	}

	m := urlMatches[0]
	if len(m) < 3 {
		return Mapping{}, errors.Errorf("failed to create mapping from %s", url)
	}
	mapping.Schema = m[1]
	mapping.Host = m[2]

	if len(m) >= 4 {
		mapping.Path = m[3]
	}

	return mapping, nil
}

func (m Mapping) HasParam(name string) bool {
	_, found := m.PathParamsSet[name]
	return found
}
