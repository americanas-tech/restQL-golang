package domain

type ResourceId string

func newResourceId(statement Statement) ResourceId {
	if statement.Alias != "" {
		return ResourceId(statement.Alias)
	}

	return ResourceId(statement.Resource)
}

type Resources map[ResourceId]interface{}

func NewResources(statements []Statement) Resources {
	resources := make(map[ResourceId]interface{})
	for _, stmt := range statements {
		index := newResourceId(stmt)
		resources[index] = stmt
	}

	return resources
}

type Debugging struct {
	Url             string
	RequestHeaders  map[string]string
	ResponseHeaders map[string]string
	Params          map[string]interface{}
	ResponseTime    int64
}

type Details struct {
	Status       int
	Success      bool
	IgnoreErrors bool
	Debug        *Debugging
}

type DoneResource struct {
	Details Details
	Result  interface{}
}

type DoneResources []interface{}
