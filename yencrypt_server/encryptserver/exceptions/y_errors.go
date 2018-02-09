package exceptions

type YError interface {
	Error()
}

type ValidationError struct {
	ErrorString string
}

func (err *ValidationError) Error() string {
	return err.ErrorString
}
