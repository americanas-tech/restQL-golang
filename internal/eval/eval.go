package eval

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
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

// ValidationError is returned by Evaluator when
// the query execution request contains invalid information.
//• Namespace: is an empty string or is not present
//• QueryRevisions name: is an empty string or is not present
//• Revision: is not a positive integer
//• Tenant: is an empty string or is not present
type ValidationError struct {
	Err error
}

func (ve ValidationError) Error() string {
	return ve.Err.Error()
}

// ParserError is returned by Evaluator when
// the asked query has invalid syntax.
type ParserError struct {
	Err error
}

func (pe ParserError) Error() string {
	return pe.Err.Error()
}

// TimeoutError is returned by Evaluator when
// the query execution time exceeds the maximum
// time defined in configuration.
type TimeoutError struct {
	Err error
}

func (te TimeoutError) Error() string {
	return te.Err.Error()
}

// MappingError is returned by Evaluator when
// the asked query references a non existing mapping.
type MappingError struct {
	Err error
}

func (me MappingError) Error() string {
	return me.Err.Error()
}
