package eval

import "github.com/b2wdigital/restQL-golang/internal/domain"

type MappingsReader interface {
	FromTenant(tenant string) (map[string]domain.Mapping, error)
}

type QueryReader interface {
	Get(namespace, id string, revision int) (string, error)
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
