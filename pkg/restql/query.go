package restql

// SavedQuery represents a query stored in database.
type SavedQuery struct {
	Text       string
	Deprecated bool
}

// QueryContext represents all data related
// to a query execution like query identification,
// input values and resource mappings.
type QueryContext struct {
	Mappings map[string]Mapping
	Options  QueryOptions
	Input    QueryInput
}

// QueryOptions represents the identity of the query being executed
type QueryOptions struct {
	Namespace string
	Id        string
	Revision  int
	Tenant    string
}

// QueryInput represents all the data
// provided by the client when requesting
// the execution of the query.
type QueryInput struct {
	Params  map[string]interface{}
	Body    interface{}
	Headers map[string]string
}
