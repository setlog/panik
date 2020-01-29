package panik

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"sync"
)

// Cause is an empty struct which signals that its position in a variadic argument list matches with a "%w" error-formatting directive.
type Cause struct{}

var describedErrors *sync.Map = &sync.Map{}

// Describe adds additional information to an ongoing panic. Subsequent calls are ineffective until Described()
// has been called.
func Describe(format string, args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	if _, isAlreadyDescribed := describedErrors.Load(r); isAlreadyDescribed {
		panic(r)
	}
	panicError := makeError(format, makeCause(r), args...)
	describedErrors.Store(panicError, nil)
	panic(panicError)
}

// Described is like Describe, but also makes the next call to Describe/Described effective again.
func Described(format string, args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	if _, isAlreadyDescribed := describedErrors.Load(r); isAlreadyDescribed {
		describedErrors.Delete(r)
		panic(r)
	}
	panic(makeError(format, makeCause(r), args...))
}

// ToError recovers from any panic if *errPtr is a nil error value and sets it to a new error which describes the recovered panic
// using the provided fmt.Errorf-compliant format with the given format args. This function panics if errPtr is nil.
func ToError(errPtr *error, format string, args ...interface{}) {
	if errPtr == nil {
		panic(fmt.Errorf("errPtr was nil"))
	}
	if *errPtr != nil {
		return
	}
	r := recover()
	if r == nil {
		return
	}
	*errPtr = makeError(format, &KnownCause{makeCause(r)}, args)
}

// ToCustomError recovers from any panic if *errPtr is a nil error value and sets it to a new error using the provided function
// with the given args. The returned error should return cause when passed to errors.Unwrap(). This function panics if errPtr is nil.
func ToCustomError(errPtr *error, newErrorFunc func(cause error, args ...interface{}) error, args ...interface{}) {
	if errPtr == nil {
		panic(fmt.Errorf("errPtr was nil"))
	}
	if *errPtr != nil {
		return
	}
	r := recover()
	if r == nil {
		return
	}
	*errPtr = newErrorFunc(&KnownCause{makeCause(r)}, args...)
}

// OnError panics with a new error if err is not nil, using given format and args in the style of fmt.Errorf.
func OnError(err error, format string, args ...interface{}) {
	if err != nil {
		panic(makeError(format, &KnownCause{cause: err}, args...))
	}
}

// Panic panics with a value of type *panik.Known which wraps an error with the provided formatted message.
func Panic(format string, args ...interface{}) {
	panic(&KnownCause{cause: fmt.Errorf(format, args...)})
}

// IsKnownCause returns true when err or one of its wrapped errors is of type *panik.Known, i.e. err came from a panic
// triggered using panik, or from To(Custom)Error (and not from somewhere else).
func IsKnownCause(err error) bool {
	var known *KnownCause
	return errors.As(err, &known)
}

// Handle handles a panic if and only if the recovered value is an error err which is or wraps an error of type *panik.Known,
// calling your provided function with the panic value type-asserted to an error value.
func Handle(handler func(r error)) {
	r := recover()
	if r == nil {
		return
	}
	if err, isError := r.(error); isError {
		var known *KnownCause
		if errors.As(err, &known) {
			handler(err)
		} else {
			panic(r)
		}
	}
	panic(r)
}

// WriteTrace recovers from any panic and writes it to the given writer, the same way that Go itself does when a goroutine
// terminates due to not having recovered from a panic, but with excessive descends into panic.go and panik.go removed.
// If there is no panic, WriteTrace does nothing.
func WriteTrace(w io.Writer) {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	w.Write([]byte(fmt.Sprintf("panic: %v\n\n%s", r, string(sb.Bytes()))))
}
