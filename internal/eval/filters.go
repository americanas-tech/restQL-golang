package eval

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/pkg/restql"
	"github.com/pkg/errors"
)

func ApplyFilters(log restql.Logger, query domain.Query, resources domain.Resources) (domain.Resources, error) {
	result := make(domain.Resources)

	for _, stmt := range query.Statements {
		resourceId := domain.NewResourceId(stmt)
		dr := resources[resourceId]

		filtered, err := applyOnlyFilters(stmt.Only, dr)
		if err != nil {
			log.Error("failed to apply filter on statement", err, "statement", fmt.Sprintf("%+#v", stmt), "done-resource", fmt.Sprintf("%+#v", dr))
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
		result := extractWithFilters(buildFilterTree(filters), resourceResult.ResponseBody)
		resourceResult.ResponseBody = result

		return resourceResult, nil
	case domain.DoneResources:
		list := make(domain.DoneResources, len(resourceResult))
		for i, r := range resourceResult {
			list[i], _ = applyOnlyFilters(filters, r)
		}
		return list, nil
	default:
		return resourceResult, errors.Errorf("resource result has unknown type %T with value: %v", resourceResult, resourceResult)
	}
}

func extractWithFilters(filters map[string]interface{}, resourceResult interface{}) interface{} {
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
			value, found := resourceResult[key]
			if !found {
				continue
			}

			if matchFilter, ok := subFilter.(domain.Match); ok {
				applyMatchFilter(matchFilter, key, value, node)
			} else if subFilter == nil {
				node[key] = value
			} else {
				subFilter, _ := subFilter.(map[string]interface{})
				node[key] = extractWithFilters(subFilter, value)
			}

		}

		return node
	case []interface{}:
		var node []interface{}
		if hasSelectAll {
			node = resourceResult
		} else {
			node = make([]interface{}, len(resourceResult))
		}

		for i, r := range resourceResult {
			node[i] = extractWithFilters(filters, r)
		}

		return node
	default:
		return resourceResult
	}
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

func applyMatchFilter(filter domain.Match, key string, value interface{}, node map[string]interface{}) {
	switch value := value.(type) {
	case []interface{}:
		var list []interface{}

		for _, v := range value {
			strVal := fmt.Sprintf("%v", v)
			match := filter.Arg.MatchString(strVal)
			if match {
				list = append(list, v)
			}
		}

		if len(list) > 0 {
			node[key] = list
		}

		return
	default:
		strVal := fmt.Sprintf("%v", value)
		match := filter.Arg.MatchString(strVal)

		if match {
			node[key] = value
		} else {
			delete(node, key)
		}

		return
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
		leaf = nil
	case domain.Match:
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
	case domain.Match:
		items, ok := s.Target().([]string)
		if !ok {
			return nil
		}

		result := make([]interface{}, len(items))
		for i, item := range items {
			if i == len(items)-1 {
				result[i] = domain.Match{Value: []string{item}, Arg: s.Arg}
			} else {
				result[i] = item
			}
		}
		return result
	default:
		return nil
	}
}

func ApplyHidden(query domain.Query, resources domain.Resources) domain.Resources {
	result := make(domain.Resources)

	for _, stmt := range query.Statements {
		if stmt.Hidden {
			continue
		}
		resourceId := domain.NewResourceId(stmt)
		dr := resources[resourceId]

		result[resourceId] = dr
	}

	return result
}
