package panik

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"sync"
)

type Cause struct{}

var describedErrors *sync.Map = &sync.Map{}

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

func ToError(errPtr *error, format string, args ...interface{}) {
	if errPtr != nil && *errPtr != nil {
		return
	}
	r := recover()
	if r == nil {
		return
	}
	err := makeError(format, r, args)
	if errPtr != nil {
		*errPtr = err
	} else {
		panic(fmt.Errorf("errPtr was nil. error was: %w", err))
	}
}

func Handle(f func(r interface{})) {
	r := recover()
	describedErrors.Delete(r)
	if r == nil {
		return
	}
	f(r)
}

// PrintStackTrace recovers from any panic and writes it to stderr, the same way that Go itself does when a
// goroutine terminates due to not having recovered from a panic, with excessive descends into panic.go and panik.go removed.
func PrintStackTrace() {
	r := recover()
	if r == nil {
		return
	}
	sb := bytes.NewBuffer(nil)
	tc := &traceCleaner{destination: sb}
	tc.Write(debug.Stack())
	os.Stderr.Write([]byte(fmt.Sprintf("panic: %v\n\n%s", r, sb.String())))
}

// ConsumeToStdLog recovers from any panic and writes it to log.Writer().
func ConsumeToStdLog() {
	ConsumeTo(log.Writer())
}

// ConsumeTo recovers from any panic and writes it to the provided writer.
func ConsumeTo(w io.Writer) {
	r := recover()
	if r == nil {
		return
	}
	message := fmt.Sprintf("%v", r)
	if message == "" {
		return
	}
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	io.WriteString(w, message)
}
