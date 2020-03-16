package runner

import "github.com/b2wdigital/restQL-golang/internal/domain"

type QueryOptions struct {
	Namespace string
	Id        string
	Revision  int
	Tenant    string
}

type QueryInput struct {
	Params  map[string]interface{}
	Headers map[string]string
}

type QueryContext struct {
	Mappings map[string]domain.Mapping
	Options  QueryOptions
	Input    QueryInput
}

type ResourceId string

func newResourceId(statement domain.Statement) ResourceId {
	if statement.Alias != "" {
		return ResourceId(statement.Alias)
	}

	return ResourceId(statement.Resource)
}

type Resources map[ResourceId]interface{}

func NewResources(statements []domain.Statement) Resources {
	resources := make(map[ResourceId]interface{})
	for _, stmt := range statements {
		index := newResourceId(stmt)
		resources[index] = stmt
	}

	return resources
}
