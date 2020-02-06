package panik

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
)

// Panic panics with an error which wraps errors.New(fmt.Sprint(args...)).
func Panic(args ...interface{}) {
	panic(&knownCause{cause: errors.New(fmt.Sprint(args...))})
}

// Panicf panics with an error which wraps fmt.Errorf(format, args...).
func Panicf(format string, args ...interface{}) {
	panic(&knownCause{cause: fmt.Errorf(format, args...)})
}

// OnError panics with err or an error which wraps err if err is not nil.
func OnError(err error) {
	if err != nil {
		if !HasKnownCause(err) {
			panic(&knownCause{cause: err})
		}
		panic(err)
	}
}

// OnErrorf panics with fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
// or an error which wraps that if err is not nil.
func OnErrorf(err error, format string, args ...interface{}) {
	if err != nil {
		if !HasKnownCause(err) {
			panic(&knownCause{cause: fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)})
		}
		panic(fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err))
	}
}

// IfError panics with panicErr or an error which wraps panicErr if err is not
// nil.
//
// You should use this function if you want to provide your own error
// implementation for panicErr.
func IfError(err error, panicErr error) {
	if err != nil {
		if !HasKnownCause(panicErr) {
			panic(&knownCause{cause: panicErr})
		}
		panic(panicErr)
	}
}

// Cause is an empty struct which signals that its position in a variadic
// argument list matches with a "%w" error-wrapping directive and that an
// error-value should take its place.
type Cause struct{}

// IfErrorf panics with an error constructed using fmt.Errorf() with the
// provided format and args or an error which wraps that if err is not nil. If
// you want to wrap err, you need to do so explicitly by matching a "%w"
// error-wrapping directive in format with a panik.Cause{} in args.
func IfErrorf(err error, format string, args ...interface{}) {
	if err != nil {
		panicErr := makeError(format, err, args...)
		if !HasKnownCause(panicErr) {
			panic(&knownCause{cause: panicErr})
		}
		panic(panicErr)
	}
}

// Wrap wraps an ongoing panic's value with an error with fmt.Sprint(args...) as
// its message.
func Wrap(args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	panic(fmt.Errorf("%s: %w", fmt.Sprint(args...), makeCause(r)))
}

// Wrapf wraps an ongoing panic's value with an error with the provided
// formatted message.
func Wrapf(format string, args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	panic(fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), makeCause(r)))
}

// ToError recovers from any panic initiated by using panik. This function
// panics if errPtr is nil.
func ToError(errPtr *error) {
	if errPtr == nil {
		panic("errPtr was nil")
	}
	if *errPtr != nil {
		return
	}
	r := recover()
	if r == nil {
		return
	}
	if !HasKnownCause(r) {
		panic(r)
	}
	*errPtr = r.(error)
}

// Recover recovers any panic and calls your provided function with the
// recovered value as parameter if it is not nil, and then panics again if the
// recovered value did not originate from one of panik's functions.
func Recover(handler func(r interface{})) {
	r := recover()
	if r == nil {
		return
	}
	handler(r)
	if !HasKnownCause(r) {
		panic(r)
	}
}

// HasKnownCause returns true when r is or wraps an error which originated from
// panik.
func HasKnownCause(r interface{}) bool {
	if err, isError := r.(error); isError {
		var known *knownCause
		return errors.As(err, &known)
	}
	return false
}

// RecoverTraceTo recovers from any panic and writes it to the given writer, the
// same way that Go itself does when a goroutine terminates due to not having
// recovered from a panic, but with excessive descends into panic.go and
// panik.go removed. If there is no panic or the panic is nil, RecoverTraceTo
// does nothing.
func RecoverTraceTo(w io.Writer) {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	w.Write([]byte(fmt.Sprintf("tracing: %v:\n%s\n", r, string(sb.Bytes()))))
}

// RecoverTrace recovers from any panic and writes it to os.Stderr, the same way
// that Go itself does when a goroutine terminates due to not having recovered
// from a panic, but with excessive descends into panic.go and panik.go removed.
// If there is no panic or the panic is nil, RecoverTrace does nothing.
func RecoverTrace() {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	os.Stderr.Write([]byte(fmt.Sprintf("tracing: %v:\n%s\n", r, string(sb.Bytes()))))
}

// ExitTraceTo recovers from any panic and writes it to the given writer, the
// same way that Go itself does when a goroutine terminates due to not having
// recovered from a panic, but with excessive descends into panic.go and
// panik.go removed, and then calls os.Exit(2). If there is no panic or the
// panic is nil, ExitTraceTo does nothing.
func ExitTraceTo(w io.Writer) {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	w.Write([]byte(fmt.Sprintf("tracing: %v:\n%s\n", r, string(sb.Bytes()))))
	os.Exit(2)
}

// ExitTrace recovers from any panic and writes it to os.Stderr, the same way
// that Go itself does when a goroutine terminates due to not having recovered
// from a panic, but with excessive descends into panic.go and panik.go removed,
// and then calls os.Exit(2). If there is no panic or the panic is nil,
// ExitTrace does nothing.
func ExitTrace() {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	os.Stderr.Write([]byte(fmt.Sprintf("tracing: %v:\n%s\n", r, string(sb.Bytes()))))
	os.Exit(2)
}
