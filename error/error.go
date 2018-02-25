package error

type ErrorString struct {
	s string
}

func (e *ErrorString) Error() string {
	return e.s
}
