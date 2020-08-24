package restql

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

var pathParamRegex = regexp.MustCompile(":([^/]+)/?")
var urlRegex = regexp.MustCompile("(https?)://([^/]+)([^?]*)\\??(.*)")

type Mapping struct {
	resourceName  string
	schema        string
	host          string
	path          string
	query         map[string]interface{}
	pathParams    []string
	pathParamsSet map[string]struct{}
}

func NewMapping(resource, url string) (Mapping, error) {
	mapping := Mapping{resourceName: resource}

	urlMatches := urlRegex.FindAllStringSubmatch(url, -1)
	if len(urlMatches) == 0 {
		return Mapping{}, errors.Errorf("failed to create mapping from %s", url)
	}

	m := urlMatches[0]
	if len(m) < 3 {
		return Mapping{}, errors.Errorf("failed to create mapping from %s", url)
	}
	mapping.schema = m[1]
	mapping.host = m[2]

	if len(m) >= 4 {
		mapping.path = m[3]
	}

	if len(m) >= 5 {
		mapping.query = parseQueryParametersInUrl(m[4])
	}

	paramsMatches := pathParamRegex.FindAllStringSubmatch(mapping.path, -1)

	pathParamsSet := make(map[string]struct{})
	pathParams := make([]string, len(paramsMatches))
	for i, m := range paramsMatches {
		paramName := m[1]
		pathParams[i] = paramName
		pathParamsSet[paramName] = struct{}{}
	}

	mapping.pathParams = pathParams
	mapping.pathParamsSet = pathParamsSet

	return mapping, nil
}

func parseQueryParametersInUrl(queryParams string) map[string]interface{} {
	if queryParams == "" {
		return nil
	}

	m := make(map[string]interface{})

	paramNames := strings.Split(queryParams, "&")
	for _, n := range paramNames {
		name := strings.Trim(n, ":")
		if name != "" {
			m[name] = struct{}{}
		}
	}

	return m
}

func (m Mapping) ResourceName() string {
	return m.resourceName
}

func (m Mapping) IsPathParam(name string) bool {
	_, found := m.pathParamsSet[name]
	return found
}

func (m Mapping) Scheme() string {
	return m.schema
}

func (m Mapping) Host() string {
	return m.host
}

func (m Mapping) IsQueryParam(name string) bool {
	_, found := m.query[name]
	return found
}

func (m Mapping) PathWithParams(params map[string]interface{}) string {
	path := m.path
	for _, pathParam := range m.pathParams {
		pathParamValue, found := params[pathParam]
		if !found {
			pathParamValue = ""
		}

		path = strings.Replace(path, fmt.Sprintf(":%v", pathParam), fmt.Sprintf("%v", pathParamValue), 1)
	}

	return path
}

func (m Mapping) QueryWithParams(params map[string]interface{}) map[string]interface{} {
	if m.query == nil {
		return nil
	}

	queryParams := make(map[string]interface{})
	for key := range m.query {
		paramValue, found := params[key]
		if !found {
			continue
		}

		queryParams[key] = paramValue
	}

	return queryParams
}
