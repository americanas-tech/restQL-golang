package restql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var pathParamRegex = regexp.MustCompile(":([^/]+)/?")
var urlRegex = regexp.MustCompile(`(https?)://([^/]+)([^?]*)\??(.*)`)

// Mapping represents the association of a name to a REST resource url.
// It support special syntax in the URL to provide dynamic value substitution, like:
//• Path parameters: can be defined by placing a colon (:) before an identifier in the URL path,
// for example "http://some.api/:id", will replace ":id" by the value of the "id" parameter
// in the query definition.
//• QueryRevisions parameters: can be defined by placing a colon (:) before an identifier in the URL query,
// for example "http://some.api?:page", will replace ":page" by the value of the "page" parameter
// in the query definition creating the URL "http://some.api?page=<value>".
type Mapping struct {
	resourceName  string
	url           string
	schema        string
	host          string
	path          string
	query         map[string]interface{}
	pathParams    []string
	pathParamsSet map[string]struct{}

	Source Source
}

// NewMapping constructs a Mapping value from a resource name
// and a canonical URL with optional identifiers for
// path and query parameters.
func NewMapping(resource, url string) (Mapping, error) {
	mapping := Mapping{resourceName: resource, url: url}

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
		mapping.query = parseQueryParametersInURL(m[4])
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

func parseQueryParametersInURL(queryParams string) map[string]interface{} {
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

// URL returns the original resource location provided
func (m Mapping) URL() string {
	return m.url
}

// ResourceName return the name associated with the resource URL
func (m Mapping) ResourceName() string {
	return m.resourceName
}

// IsPathParam returns true if the given name is a path parameter identifier
func (m Mapping) IsPathParam(name string) bool {
	_, found := m.pathParamsSet[name]
	return found
}

// Schema returns the resource URL schema
func (m Mapping) Schema() string {
	return m.schema
}

// Host returns the resource URL host
func (m Mapping) Host() string {
	return m.host
}

// IsQueryParam returns true if the given name is a query parameter identifier
func (m Mapping) IsQueryParam(name string) bool {
	_, found := m.query[name]
	return found
}

// PathWithParams takes a map of key/value pairs and use it as value source
// to replace path parameters defined by identifiers.

// PathWithParams takes a map of key/value pairs and use it to
// build a path string with the all identifiable parameters resolved.
// For example, if the mapping URL is defined as "http://some.api/:id/",
// then this method will lookup for a "id" key in the given map and use its
// value to build the result "http://some.api/<value>/".
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

// QueryWithParams takes a map of key/value pairs and use it to
// build a query parameters map with the all identifiable parameters resolved.
// For example, if the mapping URL is defined as "http://some.api?:page",
// then this method will lookup for a "page" key in the given map and use its
// value to build the result map[string]interface{}{"page": <value>}.
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
