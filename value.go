package panik

import "fmt"

// Value wraps non-error values provided to panic()
type Value struct {
	value interface{}
}

func (e *Value) Error() string {
	return fmt.Sprintf("%v", e.value)
}

func (e *Value) String() string {
	return e.Error()
}

// Unwrap returns the result of Value() if it is an error; nil otherwise.
func (e *Value) Unwrap() error {
	if err, isError := e.value.(error); isError {
		return err
	}
	return nil
}

// Value returns the value wrapped by this *Error, which has been provided to panic().
func (e *Value) Value() interface{} {
	return e.value
}
