package eval

import (
	"context"
	"errors"
	"github.com/b2wdigital/restQL-golang/v5/pkg/restql"
)

// MappingsReader is an interface implemented by types that
// can fetch a collection of mapping for the given tenant.
type MappingsReader interface {
	FromTenant(ctx context.Context, tenant string) (map[string]restql.Mapping, error)
}

// QueryReader is an interface implemented by types that
// can fetch a query for the given identification (namespace, id, revision).
type QueryReader interface {
	Get(ctx context.Context, namespace, id string, revision int) (restql.SavedQuery, error)
}

// ErrValidation is returned by Evaluator when
// the query execution request contains invalid information.
//• Namespace: is an empty string or is not present
//• QueryRevisions name: is an empty string or is not present
//• Revision: is not a positive integer
//• Tenant: is an empty string or is not present
var ErrValidation = errors.New("validation error")

// ErrParser is returned by Evaluator when
// the asked query has invalid syntax.
var ErrParser = errors.New("parsing error")

// ErrTimeout is returned by Evaluator when
// the query execution time exceeds the maximum
// time defined in configuration.
var ErrTimeout = errors.New("timeout")

// ErrMapping is returned by Evaluator when
// the asked query references a non existing mapping.
var ErrMapping = errors.New("unknown mappings")
