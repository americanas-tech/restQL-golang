package eval

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

type QueryOptions struct {
	Namespace string
	Id        string
	Revision  int
}

type QueryInput struct {
	Params  map[string]interface{}
	Headers map[string]string
}
