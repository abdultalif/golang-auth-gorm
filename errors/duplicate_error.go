package errors

type DuplicateError struct {
	Error string `json:"error"`
}

func NewDuplicateError(error string) DuplicateError {
	return DuplicateError{Error: error}
}
