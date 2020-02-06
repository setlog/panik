package panik

import (
	"fmt"
)

func makeError(format string, cause error, args ...interface{}) error {
	l := len(args)
	for i := 0; i < l; i++ {
		if _, isCause := args[i].(Cause); isCause {
			args[i] = cause
			break
		}
	}
	return fmt.Errorf(format, args...)
}

func makeCause(panicValue interface{}) error {
	if err, isError := panicValue.(error); isError {
		return err
	}
	return &Value{value: panicValue}
}

func containsCause(args ...interface{}) bool {
	for _, arg := range args {
		if _, isCause := arg.(Cause); isCause {
			return true
		}
	}
	return false
}
