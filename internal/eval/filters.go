package eval

import "github.com/b2wdigital/restQL-golang/internal/domain"

func ApplyFilters(query domain.Query, resources domain.Resources) domain.Resources {
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
