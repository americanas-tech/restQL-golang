package web

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Error is used to pass an error during the request through the
// application with web specific context.
type Error struct {
	Err    error
	Status int
}

// NewRequestError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewRequestError(err error, status int) error {
	return &Error{err, status}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (err *Error) Error() string {
	return err.Err.Error()
}
