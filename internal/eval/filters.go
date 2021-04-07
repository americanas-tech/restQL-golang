package eval

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/b2wdigital/restQL-golang/v6/internal/domain"
	"github.com/b2wdigital/restQL-golang/v6/pkg/restql"
	"github.com/pkg/errors"
)

const eot = "end-of-tree"

// ApplyFilters returns a version of the already resolved Resources
// only with the fields defined by the `only` clause.
func ApplyFilters(log restql.Logger, query domain.Query, resources domain.Resources) (domain.Resources, error) {
	result := make(domain.Resources)

	for _, stmt := range query.Statements {
		resourceID := domain.NewResourceID(stmt)
		dr := resources[resourceID]

		filtered, err := applyOnlyFilters(stmt.Only, dr)
		if err != nil {
			log.Error("failed to apply filter on statement", err, "statement", fmt.Sprintf("%+#v", stmt), "done-resource", fmt.Sprintf("%+#v", dr))
			return nil, err
		}

		result[resourceID] = filtered
	}

	return result, nil
}

func applyOnlyFilters(filters []interface{}, resourceResult interface{}) (interface{}, error) {
	if len(filters) == 0 {
		return resourceResult, nil
	}

	switch resourceResult := resourceResult.(type) {
	case restql.DoneResource:
		body := resourceResult.ResponseBody.Unmarshal()
		result, err := extractUsingFilters(buildFilterTree(filters), body)
		if err != nil {
			return nil, err
		}
		resourceResult.ResponseBody.SetValue(result)

		return resourceResult, nil
	case restql.DoneResources:
		list := make(restql.DoneResources, len(resourceResult))
		for i, r := range resourceResult {
			list[i], _ = applyOnlyFilters(filters, r)
		}
		return list, nil
	default:
		return resourceResult, errors.Errorf("resource result has unknown type %T with value: %v", resourceResult, resourceResult)
	}
}

func extractUsingFilters(filters map[string]interface{}, resourceResult interface{}) (interface{}, error) {
	filters, hasSelectAll := removeSelectAllFilter(filters)

	switch resourceResult := resourceResult.(type) {
	case map[string]interface{}:
		node := makeMapNode(hasSelectAll, resourceResult)

		for key, subFilter := range filters {
			value, found := resourceResult[key]
			if !found {
				continue
			}

			if subFilter == eot {
				node[key] = value
				continue
			}

			switch subFilter := subFilter.(type) {
			case domain.Match:
				err := applyMatchFilter(subFilter, key, value, node)
				if err != nil {
					return nil, err
				}
			case domain.FilterByRegex:
				err := applyFilterByRegex(subFilter, key, value, node)
				if err != nil {
					return nil, err
				}
			case map[string]interface{}:
				f, err := extractUsingFilters(subFilter, value)
				if err != nil {
					return nil, err
				}
				node[key] = f
			}

		}

		return node, nil
	case []interface{}:
		node := makeListNode(hasSelectAll, resourceResult)

		for i, r := range resourceResult {
			f, err := extractUsingFilters(filters, r)
			if err != nil {
				return nil, err
			}
			node[i] = f
		}

		return node, nil
	default:
		return resourceResult, nil
	}
}

func makeMapNode(hasSelectAll bool, resourceResult map[string]interface{}) map[string]interface{} {
	var node map[string]interface{}
	if hasSelectAll {
		node = resourceResult
	} else {
		node = make(map[string]interface{})
	}
	return node
}

func makeListNode(hasSelectAll bool, resourceResult []interface{}) []interface{} {
	var node []interface{}
	if hasSelectAll {
		node = resourceResult
	} else {
		node = make([]interface{}, len(resourceResult))
	}
	return node
}

func removeSelectAllFilter(filters map[string]interface{}) (map[string]interface{}, bool) {
	m := make(map[string]interface{})
	has := false

	for k, v := range filters {
		if k != "*" {
			m[k] = v
		} else {
			has = true
		}
	}

	return m, has
}

func applyFilterByRegex(fn domain.FilterByRegex, key string, value interface{}, node map[string]interface{}) error {
	var listValue []interface{}
	if valueAsList, ok := value.([]interface{}); ok {
		listValue = valueAsList
	} else {
		node[key] = value
		return nil
	}

	regex, err := parseRegex(fn.Argument(domain.FilterByRegexArgRegex).Value)
	switch {
	case err == errUnknownRegexType:
		node[key] = value
		return nil
	case err != nil:
		return err
	}

	rawPath, ok := fn.Argument(domain.FilterByRegexArgPath).Value.(string)
	if !ok {
		node[key] = value
		return nil
	}

	path := strings.Split(rawPath, ".")

	var result []interface{}
	for _, v := range listValue {
		target, found := extractValueOnPath(v, path)
		if !found {
			continue
		}

		str, err := stringify(target)
		if err != nil {
			continue
		}

		if !regex.MatchString(str) {
			continue
		}

		result = append(result, v)
	}

	if len(result) == 0 {
		node[key] = []interface{}{}
	} else {
		node[key] = result
	}

	return nil
}

func extractValueOnPath(value interface{}, path []string) (interface{}, bool) {
	if value == nil {
		return nil, false
	}

	if len(path) == 0 {
		return value, true
	}

	key := path[0]

	switch value := value.(type) {
	case map[string]interface{}:
		return extractValueOnPath(value[key], path[1:])
	default:
		return nil, false
	}
}

func applyMatchFilter(filter domain.Match, key string, value interface{}, node map[string]interface{}) error {
	matchRegex, err := parseRegex(filter.Argument(domain.MatchArgRegex).Value)
	if err != nil {
		return err
	}

	switch value := value.(type) {
	case []interface{}:
		var list []interface{}

		for _, v := range value {
			strVal, err := stringify(v)
			if err != nil {
				strVal = fmt.Sprintf("%v", v)
			}
			match := matchRegex.MatchString(strVal)
			if match {
				list = append(list, v)
			}
		}

		if len(list) > 0 {
			node[key] = list
		}

		return nil
	default:
		strVal, err := stringify(value)
		if err != nil {
			strVal = fmt.Sprintf("%v", value)
		}
		match := matchRegex.MatchString(strVal)

		if match {
			node[key] = value
		} else {
			delete(node, key)
		}

		return nil
	}
}

var errUnknownRegexType = errors.New("failed to parse match argument : unknown regex argument type")

func parseRegex(regex interface{}) (*regexp.Regexp, error) {
	switch arg := regex.(type) {
	case *regexp.Regexp:
		return arg, nil
	case string:
		return regexp.Compile(arg)
	default:
		return nil, errUnknownRegexType
	}
}

func buildFilterTree(filters []interface{}) map[string]interface{} {
	tree := make(map[string]interface{})

	for _, f := range filters {
		path := parsePath(f)
		buildPathInTree(path, tree)
	}

	return tree
}

func buildPathInTree(path []interface{}, tree map[string]interface{}) {
	if len(path) == 0 {
		return
	}

	var field string
	var leaf interface{}

	switch f := path[0].(type) {
	case string:
		field = f
		leaf = eot
	case domain.Function:
		fields, ok := f.Target().([]string)
		if !ok {
			return
		}

		field = fields[0]
		leaf = f
	}

	if len(path) == 1 {
		tree[field] = leaf
		return
	}

	if subNode, found := tree[field]; found {
		subNode, ok := subNode.(map[string]interface{})
		if !ok {
			subNode = make(map[string]interface{})
			tree[field] = subNode
		}

		buildPathInTree(path[1:], subNode)
	} else {
		subNode := make(map[string]interface{})
		tree[field] = subNode
		buildPathInTree(path[1:], subNode)
	}

}

func parsePath(s interface{}) []interface{} {
	switch s := s.(type) {
	case []string:
		items := s

		result := make([]interface{}, len(items))
		for i, item := range items {
			result[i] = item
		}
		return result
	case domain.Function:
		items, ok := s.Target().([]string)
		if !ok {
			return nil
		}

		result := make([]interface{}, len(items))
		for i, item := range items {
			if i == len(items)-1 {
				result[i] = s.Map(func(_ interface{}) interface{} {
					return []string{item}
				})
			} else {
				result[i] = item
			}
		}
		return result
	default:
		return nil
	}
}

func stringify(value interface{}) (string, error) {
	switch value := value.(type) {
	case string:
		return value, nil
	case int:
		return strconv.Itoa(value), nil
	case float64:
		return fmt.Sprintf("%.2f", value), nil
	case bool:
		return strconv.FormatBool(value), nil
	case map[string]interface{}:
		b, err := json.Marshal(value)
		return string(b), err
	case []interface{}:
		b, err := json.Marshal(value)
		return string(b), err
	default:
		return fmt.Sprintf("%v", value), nil
	}
}

// ApplyHidden returns a version of the already resolved Resources
// removing the statement results with the `hidden` clause.
func ApplyHidden(query domain.Query, resources domain.Resources) domain.Resources {
	result := make(domain.Resources)

	for _, stmt := range query.Statements {
		if stmt.Hidden {
			continue
		}
		resourceID := domain.NewResourceID(stmt)
		dr := resources[resourceID]

		result[resourceID] = dr
	}

	return result
}
