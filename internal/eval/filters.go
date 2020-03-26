package eval

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

func ApplyFilters(query domain.Query, resources domain.Resources) (domain.Resources, error) {
	result := make(domain.Resources)

	for _, stmt := range query.Statements {
		if stmt.Hidden {
			continue
		}
		resourceId := domain.NewResourceId(stmt)
		dr := resources[resourceId]

		filtered, err := applyOnlyFilters(stmt.Only, dr)
		if err != nil {
			return nil, err
		}

		result[resourceId] = filtered
	}

	return result, nil
}

func applyOnlyFilters(filters []interface{}, resourceResult interface{}) (interface{}, error) {
	if len(filters) == 0 {
		return resourceResult, nil
	}

	switch resourceResult := resourceResult.(type) {
	case domain.DoneResource:
		result, err := extractWithFilters(buildFilterTree(filters), resourceResult.Result)
		if err != nil {
			return nil, err
		}

		return domain.DoneResource{
			Details: resourceResult.Details,
			Result:  result,
		}, nil
	case domain.DoneResources:
		list := make(domain.DoneResources, len(resourceResult))
		for i, r := range resourceResult {
			list[i], _ = applyOnlyFilters(filters, r)
		}
		return list, nil
	default:
		return nil, nil
	}
}

func extractWithFilters(filters map[string]interface{}, resourceResult interface{}) (interface{}, error) {
	filters, hasSelectAll := extractSelectAllFilter(filters)

	switch resourceResult := resourceResult.(type) {
	case map[string]interface{}:
		var node map[string]interface{}
		if hasSelectAll {
			node = resourceResult
		} else {
			node = make(map[string]interface{})
		}

		for key, subFilter := range filters {
			value := resourceResult[key]

			if matchFilter, ok := subFilter.(domain.Match); ok {
				err := applyMatchFilter(matchFilter, key, value, node)
				if err != nil {
					return nil, err
				}
			} else if subFilter == nil {
				node[key] = value
			} else {
				subFilter, _ := subFilter.(map[string]interface{})
				filtered, err := extractWithFilters(subFilter, value)
				node[key] = filtered
				if err != nil {
					return nil, err
				}
			}

		}

		return node, nil
	case []interface{}:
		var node []interface{}
		if hasSelectAll {
			node = resourceResult
		} else {
			node = make([]interface{}, len(resourceResult))
		}

		for i, r := range resourceResult {
			filtered, err := extractWithFilters(filters, r)
			node[i] = filtered
			if err != nil {
				return nil, err
			}
		}

		return node, nil
	}

	return nil, nil
}

func extractSelectAllFilter(filters map[string]interface{}) (map[string]interface{}, bool) {
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

func applyMatchFilter(filter domain.Match, key string, value interface{}, node map[string]interface{}) error {
	switch value := value.(type) {
	case []interface{}:
		var list []interface{}
		regex, err := regexp.Compile(filter.Arg)
		if err != nil {
			return errors.Wrap(err, "regex matching failed")
		}

		for _, v := range value {
			strVal := fmt.Sprintf("%v", v)
			match := regex.MatchString(strVal)
			if match {
				list = append(list, v)
			}
		}

		if len(list) > 0 {
			node[key] = list
		}

		return nil
	default:
		strVal := fmt.Sprintf("%v", value)
		match, err := regexp.MatchString(filter.Arg, strVal)
		if err != nil {
			return errors.Wrap(err, "regex matching failed")
		}

		if match {
			node[key] = value
		} else {
			delete(node, key)
		}

		return nil
	}
}

func buildFilterTree(filters []interface{}) map[string]interface{} {
	tree := make(map[string]interface{})

	for _, f := range filters {
		switch f := f.(type) {
		case string:
			path := parsePath(f)
			buildPathInTree(path, tree)
		case domain.Match:
			path := parsePath(f)
			buildPathInTree(path, tree)
		}
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
		leaf = nil
	case domain.Match:
		field = f.Target
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
	case string:
		items := strings.Split(s, ".")

		result := make([]interface{}, len(items))
		for i, item := range items {
			result[i] = item
		}
		return result
	case domain.Match:
		items := strings.Split(s.Target, ".")

		result := make([]interface{}, len(items))
		for i, item := range items {
			if i == len(items)-1 {
				result[i] = domain.Match{Target: item, Arg: s.Arg}
			} else {
				result[i] = item
			}
		}
		return result
	default:
		return nil
	}
}
