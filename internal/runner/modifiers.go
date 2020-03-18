package runner

import "github.com/b2wdigital/restQL-golang/internal/domain"

func ApplyModifiers(modifiers domain.Modifiers, resources Resources) Resources {
	for resourceId, stmt := range resources {
		if stmt, ok := stmt.(domain.Statement); ok {
			stmt.CacheControl = applyCacheModifiers(modifiers, stmt)
			stmt.Timeout = applyTimeoutModifier(modifiers, stmt)
			resources[resourceId] = stmt
		}
	}

	return resources
}

func applyTimeoutModifier(modifiers domain.Modifiers, statement domain.Statement) interface{} {
	timeout := modifiers["timeout"]
	currentTimeout := statement.Timeout

	if currentTimeout == nil {
		return timeout
	}

	return currentTimeout
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
