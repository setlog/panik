package panik

import "fmt"

import "strings"

// Error wraps non-error values provided to panic()
type Error struct {
	value interface{}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v", e.value)
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

func makeError(format string, panicValue interface{}, args ...interface{}) error {
	panicError := makeCause(panicValue)
	l := len(args)
	for i := 0; i < l; i++ {
		if _, isCause := args[i].(Cause); isCause {
			args[i] = panicError
		}
	}
	if !hasErrorFormattingDirective(format) {
		format += ": %w"
		args = append(args, panicError)
	}
	return fmt.Errorf(format, args...)
}

func makeCause(panicValue interface{}) error {
	if err, isError := panicValue.(error); isError {
		return err
	}
	return &Error{value: panicValue}
}

func hasErrorFormattingDirective(format string) bool {
	for {
		i := strings.Index(format, "%w")
		if i == -1 {
			return false
		}
		if i == 0 || format[i-1] != '%' {
			return true
		}
		format = format[i+2:]
	}
}
