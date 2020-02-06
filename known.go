package panik

// knownCause is an error-wrapper which signals that it originates from a call
// to one of panik's functions.
type knownCause struct {
	cause error
}

func (e *knownCause) Error() string {
	return e.cause.Error()
}

func (e *knownCause) String() string {
	return e.Error()
}

func (e *knownCause) Unwrap() error {
	return e.cause
}
