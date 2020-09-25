package runner

import (
	"github.com/b2wdigital/restQL-golang/v4/internal/domain"
)

// ApplyModifiers transforms an unresolved Resources collection and
// set the query level cache directives into each statement without
// cache directives.
func ApplyModifiers(resources domain.Resources, modifiers domain.Modifiers) domain.Resources {
	for resourceID, stmt := range resources {
		if stmt, ok := stmt.(domain.Statement); ok {
			stmt.CacheControl = applyCacheModifiers(modifiers, stmt)
			resources[resourceID] = stmt
		}
	}

	return resources
}

func applyCacheModifiers(modifiers domain.Modifiers, statement domain.Statement) domain.CacheControl {
	cc := statement.CacheControl

	cacheControl, found := modifiers["cache-control"]
	if cc.MaxAge == nil && found {
		cc.MaxAge = cacheControl
	}

	maxAge, found := modifiers["max-age"]
	if cc.MaxAge == nil && found {
		cc.MaxAge = maxAge
	}

	smaxAge, found := modifiers["s-max-age"]
	if cc.SMaxAge == nil && found {
		cc.SMaxAge = smaxAge
	}

	return cc
}
