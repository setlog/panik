package panik

// knownCause is an error-wrapper which signals that it originates from a call
// to one of panik's functions.
type knownCause struct {
	message            string
	cause              error
	deescalatedToError bool
}

func (e *knownCause) Error() string {
	return e.message
}

func (e *knownCause) String() string {
	return e.Error()
}

func (e *knownCause) Unwrap() error {
	return e.cause
}

func makeCause(panicValue interface{}, stripKnownCause bool) error {
	if err, isError := panicValue.(error); isError && err != nil {
		if stripKnownCause {
			if knownErr, isKnownErr := err.(*knownCause); isKnownErr {
				return knownErr.cause
			}
		}
		return err
	}
	return &Value{value: panicValue}
}
