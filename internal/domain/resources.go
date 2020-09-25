package domain

import "github.com/b2wdigital/restQL-golang/v4/pkg/restql"

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

// ResourceCacheControlValue represents the values a cache control
// directive is able to have.
// It can either be present and have a integer time value or
// not be present in the upstream response.
type ResourceCacheControlValue struct {
	Exist bool
	Time  int
}

// ResourceCacheControl represent cache control directives
// returned by upstream during statement resolution.
type ResourceCacheControl struct {
	NoCache bool
	MaxAge  ResourceCacheControlValue
	SMaxAge ResourceCacheControlValue
}

// Details represents metadata about the statement result.
type Details struct {
	Status       int
	Success      bool
	IgnoreErrors bool
	CacheControl ResourceCacheControl
	Debug        *Debugging
}

// DoneResource represents a statement result.
type DoneResource struct {
	Status          int
	Success         bool
	IgnoreErrors    bool
	CacheControl    ResourceCacheControl
	Method          string
	URL             string
	RequestParams   map[string]interface{}
	RequestHeaders  map[string]string
	RequestBody     interface{}
	ResponseHeaders map[string]string
	ResponseBody    *restql.ResponseBody
	ResponseTime    int64
}

// DoneResources represents a multiplexed statement result.
type DoneResources []interface{}
