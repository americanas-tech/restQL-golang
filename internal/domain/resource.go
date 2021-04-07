package domain

import "github.com/b2wdigital/restQL-golang/v6/pkg/restql"

// ResourceID is an unique identifier used by a statement.
// If an alias is present, it is used. Otherwise, the resource
// name.
type ResourceID string

// NewResourceID make a ResourceID from a Statement.
func NewResourceID(statement Statement) ResourceID {
	if statement.Alias != "" {
		return ResourceID(statement.Alias)
	}

	return ResourceID(statement.Resource)
}

// Resources represents the index of statements
// already resolved or to be resolved.
type Resources map[ResourceID]interface{}

// NewResources constructs a Resources collection
// from a slice of statements.
func NewResources(statements []Statement) Resources {
	resources := make(map[ResourceID]interface{})
	for _, stmt := range statements {
		index := NewResourceID(stmt)
		resources[index] = stmt
	}

	return resources
}

// Debugging represents the collection of information
// about the statement result.
// This is only used when the client enable it during
// query execution.
type Debugging struct {
	Method          string
	Url             string
	RequestHeaders  map[string]string
	ResponseHeaders map[string]string
	Params          map[string]interface{}
	RequestBody     interface{}
	ResponseTime    int64
}

// Details represents metadata about the statement result.
type Details struct {
	Status       int
	Success      bool
	IgnoreErrors bool
	CacheControl restql.ResourceCacheControl
	Debug        *Debugging
}
