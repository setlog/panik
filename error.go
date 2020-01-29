package panik

import (
	"fmt"
	"regexp"
)

func makeError(format string, cause error, args ...interface{}) error {
	l := len(args)
	for i := 0; i < l; i++ {
		if _, isCause := args[i].(Cause); isCause {
			args[i] = cause
		}
	}
	if !hasErrorFormattingDirective.MatchString(format) {
		format += ": %w"
		args = append(args, cause)
	}
	return fmt.Errorf(format, args...)
}

func makeCause(panicValue interface{}) error {
	if err, isError := panicValue.(error); isError {
		return err
	}
	return &Value{value: panicValue}
}

var hasErrorFormattingDirective *regexp.Regexp = regexp.MustCompile("(([^%]|^)(%%)*%w)")
