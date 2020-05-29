package domain

type ResourceId string

func NewResourceId(statement Statement) ResourceId {
	if statement.Alias != "" {
		return ResourceId(statement.Alias)
	}

	return ResourceId(statement.Resource)
}

type Resources map[ResourceId]interface{}

func NewResources(statements []Statement) Resources {
	resources := make(map[ResourceId]interface{})
	for _, stmt := range statements {
		index := NewResourceId(stmt)
		resources[index] = stmt
	}

	return resources
}

type Debugging struct {
	Method          string
	Url             string
	RequestHeaders  map[string]string
	ResponseHeaders map[string]string
	Params          map[string]interface{}
	RequestBody     interface{}
	ResponseTime    int64
}

type ResourceCacheControlValue struct {
	Exist bool
	Time  int
}

type ResourceCacheControl struct {
	NoCache bool
	MaxAge  ResourceCacheControlValue
	SMaxAge ResourceCacheControlValue
}

type Details struct {
	Status       int
	Success      bool
	IgnoreErrors bool
	CacheControl ResourceCacheControl
	Debug        *Debugging
}

type DoneResource struct {
	Details Details
	Result  interface{}

	Status          int
	Success         bool
	IgnoreErrors    bool
	CacheControl    ResourceCacheControl
	Method          string
	Url             string
	RequestParams   map[string]interface{}
	RequestHeaders  map[string]string
	RequestBody     interface{}
	ResponseHeaders map[string]string
	ResponseBody    interface{}
	ResponseTime    int64
}

type DoneResources []interface{}
