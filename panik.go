package panik

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
)

var describedErrors *sync.Map = &sync.Map{}

func Describe(format string, args ...interface{}) {
	r := recover()
	if r == nil {
		return
	}
	_, isAlreadyDescribed := describedErrors.Load(r)
	if isAlreadyDescribed {
		panic(r)
	}
	var panicError error
	if err, isError := r.(error); isError {
		panicError = err
	} else {
		panicError = &Error{value: r}
	}
	args = append(args, panicError)
	panicError = fmt.Errorf(format+": %w", args...)
	describedErrors.Store(panicError, nil)
	panic(panicError)
}

func Consolidate(format string, args ...interface{}) {
	defer ConsolidateAsIs()
	Describe(format, args)
}

func ConsolidateAsIs() {
	r := recover()
	describedErrors.Delete(r)
	if r == nil {
		return
	}
	panic(r)
}

func ToNewError(errPtr *error, format string, args ...interface{}) {
	if errPtr == nil {
		log.Printf("errPtr was nil")
		return
	}
	if *errPtr != nil {
		return
	}
	r := recover()
	if r == nil {
		return
	}
	var panicError error
	if err, isError := r.(error); isError {
		panicError = err
	} else {
		panicError = &Error{value: r}
	}
	args = append(args, panicError)
	*errPtr = fmt.Errorf(format, args...)
}

func AsError(errPtr *error) {
	if errPtr == nil {
		log.Printf("errPtr was nil")
		return
	}
	if *errPtr != nil {
		return
	}
	r := recover()
	if r == nil {
		return
	}
	if err, isError := r.(error); isError {
		*errPtr = err
	} else {
		*errPtr = fmt.Errorf("%v", r)
	}
}

// TODO: Provide clean stacktrace
func Handle(f func(r interface{})) {
	r := recover()
	if r == nil {
		return
	}
	describedErrors.Delete(r)
	f(r)
}

func ConsumeToStandardLogger() {
	r := recover()
	if r == nil {
		return
	}
	log.Printf("Ignoring panic: %v", r)
}

func ConsumeToWriter(w io.Writer) {
	r := recover()
	if r == nil {
		return
	}
	message := fmt.Sprintf("%v", r)
	if message == "" {
		return
	}
	if strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	io.WriteString(w, message)
}
