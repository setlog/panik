package panik

import (
	"bytes"
	"fmt"
	"io"
	"runtime/debug"
	"sync"
)

// Cause is an empty struct which signals to Describe(), Described() and ToError() which argument to replace
// with the underlying error to match the "%w" error-formatting directive.
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
	panicError := makeError(format, r, args...)
	describedErrors.Store(panicError, nil)
	panic(panicError)
}

// Described is like Describe, but also makes the next call to Describe effective again.
func Described(format string, args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	if _, isAlreadyDescribed := describedErrors.Load(r); isAlreadyDescribed {
		describedErrors.Delete(r)
		panic(r)
	}
	panic(makeError(format, r, args...))
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
	*errPtr = makeError(format, r, args)
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
	*errPtr = newErrorFunc(makeCause(r), args...)
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
