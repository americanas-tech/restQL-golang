package runner

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/v5/internal/domain"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
	"github.com/pkg/errors"
)

// ErrInvalidDependsOnTarget represents an error when a depends-on target
// references an unknown statement.
var ErrInvalidDependsOnTarget = errors.New("depends-on targets an unknown resource")

// ResolveDependsOn takes an unresolved Resource collection and
// find if a resource dependency was successfully resolved
func ResolveDependsOn(resources domain.Resources, doneResources domain.Resources) domain.Resources {
	for resourceID, stmt := range resources {
		resources[resourceID] = resolveDependsOnIntoStatement(stmt, doneResources)
	}

	return resources
}

func resolveDependsOnIntoStatement(stmt interface{}, doneResources domain.Resources) interface{} {
	switch stmt := stmt.(type) {
	case domain.Statement:
		target := stmt.DependsOn.Target
		if target == "" {
			stmt.DependsOn.Resolved = true
			return stmt
		}

		dr := doneResources[domain.ResourceID(target)]
		stmt.DependsOn.Resolved = isResolved(dr)
		return stmt
	case []interface{}:
		result := make([]interface{}, len(stmt))
		for i, s := range stmt {
			result[i] = resolveDependsOnIntoStatement(s, doneResources)
		}
		return result
	default:
		return stmt
	}
}

func isResolved(doneResource interface{}) bool {
	switch done := doneResource.(type) {
	case restql.DoneResource:
		return done.IgnoreErrors || (done.Status >= 200 && done.Status <= 399)
	case restql.DoneResources:
		resolved := false
		for _, d := range done {
			resolved = resolved || isResolved(d)
		}
		return resolved
	default:
		return false
	}
}

// ValidateDependsOnTarget returns an error if a depends-on
// target references an unknown statement.
func ValidateDependsOnTarget(resources domain.Resources) error {
	for _, stmt := range resources {
		err := validateDependsOnIntoStatement(stmt, resources)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateDependsOnIntoStatement(stmt interface{}, resources domain.Resources) error {
	switch stmt := stmt.(type) {
	case domain.Statement:
		target := stmt.DependsOn.Target
		if target == "" {
			return nil
		}

		_, found := resources[domain.ResourceID(target)]
		if !found {
			return fmt.Errorf("%w: %s", ErrInvalidDependsOnTarget, target)
		}

		return nil
	case []interface{}:
		for _, s := range stmt {
			err := validateDependsOnIntoStatement(s, resources)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return nil
}
