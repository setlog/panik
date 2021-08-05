package panik

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
)

// Panic panics with an error which wraps r if r is an error and an error
// which wraps &panik.Value{r} if r is not an error.
func Panic(r interface{}) {
	c := makeCause(r, false)
	panic(&knownCause{cause: c, message: c.Error()})
}

// Panicf panics with an error which wraps fmt.Errorf(format, args...).
func Panicf(format string, args ...interface{}) {
	c := fmt.Errorf(format, args...)
	panic(&knownCause{cause: c, message: c.Error()})
}

// OnError panics with an error which wraps panicErr and returns
// panicErr.Error() for Error() if panicErr is not nil.
func OnError(panicErr error) {
	if panicErr != nil {
		panic(&knownCause{cause: panicErr, message: panicErr.Error()})
	}
}

// OnErrore panics with an error which wraps panicErr and returns
// fmt.Sprintf("%v: %v", panicErr, err) for Error() if neither err
// nor panicErr are nil.
//
// You may use this function if you want to provide your own error
// implementation for panicErr in reaction to err.
func OnErrore(err error, panicErr error) {
	if err != nil && panicErr != nil {
		panic(&knownCause{cause: panicErr, message: fmt.Sprintf("%v: %v", panicErr, err)})
	}
}

// OnErrorfw panics with an error which wraps panicErr and returns
// fmt.Sprintf("%s: %v", fmt.Sprintf(format, args...), err) for Error()
// if panicErr is not nil.
//
// You may use this function instead of OnError if you wish to supply
// more information regarding the circumstances of the error.
func OnErrorfw(panicErr error, format string, args ...interface{}) {
	if panicErr != nil {
		panic(&knownCause{cause: panicErr, message: fmt.Sprintf("%s: %v", fmt.Sprintf(format, args...), panicErr)})
	}
}

// OnErrorfv panics with an error which does not wrap panicErr and returns
// fmt.Sprintf("%s: %v", fmt.Sprintf(format, args...), err) for Error()
// if panicErr is not nil.
//
// You may use this function if you do not wish expose the type and data
// of panicErr (apart from the string returned from Error()) to the caller.
func OnErrorfv(panicErr error, format string, args ...interface{}) {
	if panicErr != nil {
		cause := fmt.Errorf("%v", panicErr)
		panic(&knownCause{cause: cause, message: fmt.Sprintf("%s: %v", fmt.Sprintf(format, args...), panicErr)})
	}
}

// Wrap wraps an ongoing panic's value with an error with
// fmt.Sprint(args...) as its message.
//
// Wrap may only be used in a defer-statement.
func Wrap(args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	err := fmt.Errorf("%s: %w", fmt.Sprint(args...), makeCause(r, true))
	if caused(r) {
		panic(&knownCause{cause: err, message: err.Error()})
	} else {
		panic(err)
	}
}

// Wrapf wraps an ongoing panic's value with an error with
// fmt.Sprintf(format, args...) as its message.
//
// Wrapf may only be used in a defer-statement.
func Wrapf(format string, args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	err := fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), makeCause(r, true))
	if caused(r) {
		panic(&knownCause{cause: err, message: err.Error()})
	} else {
		panic(err)
	}
}

// ToError recovers from any panic which originated from panik and writes the
// recovered error to *errPtr.
//
// This function panics if errPtr is nil and does nothing if *errPtr is non-nil.
//
// ToError may only be used in a defer-statement.
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
	if !caused(r) {
		panic(r)
	}
	r.(*knownCause).deescalatedToError = true
	*errPtr = r.(*knownCause).cause
}

// ToErrorWithTrace recovers from any panic which originated from panik and writes
// an error which wraps the recovered error to *errPtr and contains the stack trace
// of the panic in its message.
//
// This function panics if errPtr is nil and does nothing if *errPtr is non-nil.
//
// ToErrorWithTrace may only be used in a defer-statement.
func ToErrorWithTrace(errPtr *error) {
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
	if !caused(r) {
		panic(r)
	}
	r.(*knownCause).deescalatedToError = true
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	*errPtr = fmt.Errorf("%w:\n%s", r.(*knownCause).cause, sb.String())
}

// Caused returns true if r is an error which originated from panik.
func Caused(r interface{}) bool {
	knownErr, isKnown := r.(*knownCause)
	return isKnown && knownErr != nil
}

func caused(r interface{}) bool {
	knownErr, isKnown := r.(*knownCause)
	return isKnown && knownErr != nil && !knownErr.deescalatedToError
}

// RecoverTraceTo recovers from any panic and writes it to the given writer, the
// same way that Go itself does when a goroutine terminates due to not having
// recovered from a panic, but with excessive descends into panic.go and panik
// removed. If there is no panic or the panic is nil, RecoverTraceTo does
// nothing.
//
// RecoverTraceTo may only be used in a defer-statement.
func RecoverTraceTo(w io.Writer) {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	w.Write([]byte(fmt.Sprintf("recovered: %v:\n%s\n", r, sb.String())))
}

// RecoverTraceToDefaultLogger is like RecoverTraceTo, but always uses log.Default().Writer()
// as the output target.
//
// RecoverTraceToDefaultLogger may only be used in a defer-statement.
func RecoverTraceToDefaultLogger() {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	log.Default().Writer().Write([]byte(fmt.Sprintf("recovered: %v:\n%s\n", r, sb.String())))
}

// RecoverTraceFunc recovers from any panic and calls provided function with a stack trace,
// formatted the same way that Go itself does when a goroutine terminates due to not having
// recovered from a panic, but with excessive descends into panic.go and panik removed. If
// there is no panic or the panic is nil, RecoverTraceFunc does nothing.
//
// RecoverTraceFunc may only be used in a defer-statement.
func RecoverTraceFunc(f func(trace string)) {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	f(fmt.Sprintf("recovered: %v:\n%s\n", r, sb.String()))
}

// ExitTraceTo is like RecoverTraceTo, but also calls os.Exit(2) after writing to w.
//
// ExitTraceTo may only be used in a defer-statement.
func ExitTraceTo(w io.Writer) {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	w.Write([]byte(fmt.Sprintf("panic: %v:\n%s\n", r, sb.String())))
	os.Exit(2)
}

// ExitTraceToDefaultLogger is like RecoverTraceToDefaultLogger, but also calls os.Exit(2) after writing to log.Default().Writer().
//
// ExitTraceToDefaultLogger may only be used in a defer-statement.
func ExitTraceToDefaultLogger() {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	log.Default().Writer().Write([]byte(fmt.Sprintf("panic: %v:\n%s\n", r, sb.String())))
	os.Exit(2)
}

// ExitTraceFunc is like RecoverTraceFunc, but also calls os.Exit(2) after returning from f.
//
// ExitTraceFunc may only be used in a defer-statement.
func ExitTraceFunc(f func(trace string)) {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	f(fmt.Sprintf("panic: %v:\n%s\n", r, sb.String()))
	os.Exit(2)
}
