package eval

import (
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"strings"
)

var Leaf interface{} = nil

func ApplyFilters(query domain.Query, resources domain.Resources) domain.Resources {
	result := make(domain.Resources)

	for _, stmt := range query.Statements {
		if stmt.Hidden {
			continue
		}
		resourceId := domain.NewResourceId(stmt)
		dr := resources[resourceId]
		result[resourceId] = applyOnlyFilters(stmt.Only, dr)
	}

	return result
}

func applyOnlyFilters(filters []interface{}, resourceResult interface{}) interface{} {
	if len(filters) == 0 {
		return resourceResult
	}

	switch resourceResult := resourceResult.(type) {
	case domain.DoneResource:
		return domain.DoneResource{
			Details: resourceResult.Details,
			Result:  extractWithFilters(buildFilterTree(filters), resourceResult.Result),
		}
	case domain.DoneResources:
		list := make(domain.DoneResources, len(resourceResult))
		for i, r := range resourceResult {
			list[i] = applyOnlyFilters(filters, r)
		}
		return list
	default:
		return nil
	}
}

func extractWithFilters(filters map[string]interface{}, resourceResult interface{}) interface{} {
	switch resourceResult := resourceResult.(type) {
	case map[string]interface{}:
		node := make(map[string]interface{})

		for key, subFilter := range filters {
			value := resourceResult[key]
			if subFilter == Leaf {
				node[key] = value
			} else {
				subFilter, _ := subFilter.(map[string]interface{})
				node[key] = extractWithFilters(subFilter, value)
			}

		}

		return node
	case []interface{}:
		node := make([]interface{}, len(resourceResult))

		for i, r := range resourceResult {
			node[i] = extractWithFilters(filters, r)
		}

		return node
	}

	return nil
}

func buildFilterTree(filters []interface{}) map[string]interface{} {
	tree := make(map[string]interface{})

	for _, f := range filters {
		switch f := f.(type) {
		case string:
			path := parsePath(f)
			buildPathInTree(path, tree)
		}
	}

	return tree
}

func buildPathInTree(path []string, tree map[string]interface{}) {
	if len(path) == 0 {
		return
	}

	field := path[0]
	if len(path) == 1 {
		tree[field] = Leaf
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

func parsePath(s string) []string {
	return strings.Split(s, ".")
}
