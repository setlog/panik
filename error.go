package panik

import "fmt"

// Error wraps non-error values provided to panic()
type Error struct {
	value interface{}
}

func (e *Error) Error() string {
	if err, isError := e.value.(error); isError && err != nil {
		return "panic value: " + err.Error()
	}
	return fmt.Sprintf("panic value: %v", e)
}

func (e *Error) String() string {
	return e.Error()
}

// Unwrap returns the result of Value() if it is an error; nil otherwise.
func (e *Error) Unwrap() error {
	if err, isError := e.value.(error); isError {
		return err
	}
	return nil
}

// Value returns the value wrapped by this *Error, which has been provided to panic().
func (e *Error) Value() interface{} {
	return e.value
}
