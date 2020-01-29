package panik

import "fmt"

// value wraps non-error values provided to panic()
type value struct {
	value interface{}
}

func (e *value) Error() string {
	return fmt.Sprintf("%v", e.value)
}

func (e *value) String() string {
	return e.Error()
}

// Unwrap returns the result of Value() if it is an error; nil otherwise.
func (e *value) Unwrap() error {
	if err, isError := e.value.(error); isError {
		return err
	}
	return nil
}

// Value returns the value wrapped by this *Error, which has been provided to panic().
func (e *value) Value() interface{} {
	return e.value
}
