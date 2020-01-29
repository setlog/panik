package panik

// KnownCause is an error-wrapper which signals that it originates from a call to To(Custom)Error, OnError or Start.
type KnownCause struct {
	cause error
}

func (e *KnownCause) Error() string {
	return e.cause.Error()
}

func (e *KnownCause) String() string {
	return e.Error()
}

func (e *KnownCause) Unwrap() error {
	return e.cause
}
