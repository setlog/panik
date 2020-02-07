package panik

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
)

// Panic panics with an error which wraps r if it is an error and &Value{r} if r
// is not an error.
func Panic(r interface{}) {
	c := makeCause(r)
	panic(&knownCause{cause: c, message: c.Error()})
}

// Panicf panics with an error which wraps fmt.Errorf(format, args...).
func Panicf(format string, args ...interface{}) {
	c := fmt.Errorf(format, args...)
	panic(&knownCause{cause: c, message: c.Error()})
}

// OnError panics with an error which wraps err if err is not nil.
func OnError(err error) {
	if err != nil {
		panic(&knownCause{cause: err, message: err.Error()})
	}
}

// OnErrore panics with an error which wraps panicErr if neither err nor
// panicErr are nil.
//
// You should use this function if you want to provide your own error
// implementation for panicErr.
func OnErrore(err error, panicErr error) {
	if err != nil && panicErr != nil {
		panic(&knownCause{cause: panicErr, message: fmt.Sprintf("%v: %v", panicErr, err)})
	}
}

// OnErrorfw panics with an error which wraps err and returns
// fmt.Sprintf("%s: %v", fmt.Sprintf(format, args...), err) for Error().
func OnErrorfw(err error, format string, args ...interface{}) {
	if err != nil {
		panic(&knownCause{cause: err, message: fmt.Sprintf("%s: %v", fmt.Sprintf(format, args...), err)})
	}
}

// OnErrorfv panics with an error which wraps fmt.Errorf("%s: %v",
// fmt.Sprintf(format, args...), err) if err is not nil.
func OnErrorfv(err error, format string, args ...interface{}) {
	if err != nil {
		panic(&knownCause{cause: nil, message: fmt.Sprintf("%s: %v", fmt.Sprintf(format, args...), err)})
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

// Wrapf wraps an ongoing panic's value with an error with fmt.Sprintf(format,
// args...) as its message.
func Wrapf(format string, args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	panic(fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), makeCause(r)))
}

// ToError recovers from any panic which originated from panik. This function
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
	if !Caused(r) {
		panic(r)
	}
	*errPtr = r.(error)
}

// Caused returns true when r is or wraps an error which originated from panik.
func Caused(r interface{}) bool {
	if err, isError := r.(error); isError {
		var known *knownCause
		return errors.As(err, &known)
	}
	return false
}

// RecoverTraceTo recovers from any panic and writes it to the given writer, the
// same way that Go itself does when a goroutine terminates due to not having
// recovered from a panic, but with excessive descends into panic.go and panik
// removed. If there is no panic or the panic is nil, RecoverTraceTo does
// nothing.
func RecoverTraceTo(w io.Writer) {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	w.Write([]byte(fmt.Sprintf("recovered: %v:\n%s\n", r, string(sb.Bytes()))))
}

// RecoverTrace recovers from any panic and writes it to os.Stderr, the same way
// that Go itself does when a goroutine terminates due to not having recovered
// from a panic, but with excessive descends into panic.go and panik removed. If
// there is no panic or the panic is nil, RecoverTrace does nothing.
func RecoverTrace() {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	os.Stderr.Write([]byte(fmt.Sprintf("recovered: %v:\n%s\n", r, string(sb.Bytes()))))
}

// ExitTraceTo recovers from any panic and writes it to the given writer, the
// same way that Go itself does when a goroutine terminates due to not having
// recovered from a panic, but with excessive descends into panic.go and panik
// removed, and then calls os.Exit(2). If there is no panic or the panic is nil,
// ExitTraceTo does nothing.
func ExitTraceTo(w io.Writer) {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	w.Write([]byte(fmt.Sprintf("fatal: %v:\n%s\n", r, string(sb.Bytes()))))
	os.Exit(2)
}

// ExitTrace recovers from any panic and writes it to os.Stderr, the same way
// that Go itself does when a goroutine terminates due to not having recovered
// from a panic, but with excessive descends into panic.go and panik removed,
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
	os.Stderr.Write([]byte(fmt.Sprintf("fatal: %v:\n%s\n", r, string(sb.Bytes()))))
	os.Exit(2)
}
