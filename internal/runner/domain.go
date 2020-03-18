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

type Debugging struct {
	Url             string                 `json:"url,omitempty"`
	RequestHeaders  map[string]string      `json:"request-headers,omitempty"`
	ResponseHeaders map[string]string      `json:"response-headers,omitempty"`
	Params          map[string]interface{} `json:"params,omitempty"`
}

type Details struct {
	Status  int        `json:"status"`
	Success bool       `json:"success"`
	Debug   *Debugging `json:"debug,omitempty"`
}

type Response struct {
	Details Details     `json:"details"`
	Result  interface{} `json:"result"`
}
