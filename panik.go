package panik

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
)

// Cause is an empty struct which signals that its position in a variadic argument list matches
// with a "%w" error-wrapping directive and that an error-value should take its place.
type Cause struct{}

// Errorf formats an ongoing panic in the style of fmt.Errorf(). If you want to preserve its
// cause, you need to use the "%w" error-wrapping directive and provide panik.Cause{} as a
// matching argument.
func Errorf(format string, args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	if HasKnownCause(r) && !(hasErrorFormattingDirective.MatchString(format) && containsCause(args...)) {
		panic(&knownCause{cause: fmt.Errorf(format, args...)})
	}
	panic(makeError(format, makeCause(r), args...))
}

// ToError recovers from any panic which is or wraps a *panik.knownCause if *errPtr is nil
// and sets it to the recovered error.
//
// This function panics if errPtr is nil.
func ToError(errPtr *error) {
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
	if !HasKnownCause(r) {
		panic(r)
	}
	*errPtr = r.(error)
}

// OnError panics with &knownCause{cause: err} if err is not nil.
func OnError(err error) {
	if err != nil {
		panic(&knownCause{cause: err})
	}
}

// IfError panics with &knownCause{cause: panicErr} if err is not nil.
func IfError(err error, panicErr error) {
	if err != nil {
		panic(&knownCause{cause: panicErr})
	}
}

// IfErrorf panics with &knownCause{cause: fmt.Errorf(format, args...)} if err is not nil.
func IfErrorf(err error, format string, args ...interface{}) {
	if err != nil {
		panic(&knownCause{cause: fmt.Errorf(format, args...)})
	}
}

// Panicf panics with a value of type *panik.knownCause which wraps a new error
// with the provided formatted message.
func Panicf(format string, args ...interface{}) {
	panic(&knownCause{cause: fmt.Errorf(format, args...)})
}

// HasKnownCause returns true when r or one of its wrapped errors is of type
// *panik.knownCause, i.e. r came from a panic triggered using panik.
func HasKnownCause(r interface{}) bool {
	if err, isError := r.(error); isError {
		var known *knownCause
		return errors.As(err, &known)
	}
	return false
}

// Handle recovers a panic if and only if the recovered value is an error which is or wraps a *panik.knownCause
// and calls your provided handler function with it as a parameter.
//
// Use cases of Handle() include:
// - Calling panik.OnError() with your own implementation of the error interface.
// - Doing cleanup/finalization based on the type of r.
//
// Clean-up of state should happen in a traditional defer func(){}()-call.
func Handle(handler func(r error)) {
	r := recover()
	if r == nil {
		return
	}
	if !HasKnownCause(r) {
		panic(r)
	}
	handler(r.(error))
}

// RecoverTraceTo recovers from any panic and writes it to the given writer, the same way that Go itself does when a goroutine
// terminates due to not having recovered from a panic, but with excessive descends into panic.go and panik.go removed.
// If there is no panic or the panic is nil, RecoverTraceTo does nothing.
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

// RecoverTrace recovers from any panic and writes it to os.Stderr, the same way that Go itself does when a goroutine
// terminates due to not having recovered from a panic, but with excessive descends into panic.go and panik.go removed.
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

// ExitTraceTo recovers from any panic and writes it to the given writer, the same way that Go itself does when a goroutine
// terminates due to not having recovered from a panic, but with excessive descends into panic.go and panik.go removed.
// If there is no panic or the panic is nil, ExitTraceTo does nothing.
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

// ExitTrace recovers from any panic and writes it to os.Stderr, the same way that Go itself does when a goroutine
// terminates due to not having recovered from a panic, but with excessive descends into panic.go and panik.go removed.
// If there is no panic or the panic is nil, ExitTrace does nothing.
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
