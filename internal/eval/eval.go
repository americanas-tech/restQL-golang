package eval

import (
	"context"
	"github.com/b2wdigital/restQL-golang/v4/pkg/restql"
)

type MappingsReader interface {
	FromTenant(ctx context.Context, tenant string) (map[string]restql.Mapping, error)
}

type QueryReader interface {
	Get(ctx context.Context, namespace, id string, revision int) (restql.SavedQuery, error)
}

type ValidationError struct {
	Err error
}

func (ve ValidationError) Error() string {
	return ve.Err.Error()
}

type NotFoundError struct {
	Err error
}

func (ne NotFoundError) Error() string {
	return ne.Err.Error()
}

type ParserError struct {
	Err error
}

func (pe ParserError) Error() string {
	return pe.Err.Error()
}

type TimeoutError struct {
	Err error
}

func (te TimeoutError) Error() string {
	return te.Err.Error()
}

type MappingError struct {
	Err error
}

func (me MappingError) Error() string {
	return me.Err.Error()
}
