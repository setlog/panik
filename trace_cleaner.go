package panik

import (
	"io"
	"regexp"
	"strings"
)

var verboseRegExps []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`^panic\(.*$`),
	regexp.MustCompile(`^runtime/debug.Stack\(.*$`),
	regexp.MustCompile(`^github.com/setlog/panik\..*\(.*$`),
}

type traceCleaner struct {
	destination io.Writer
	line        string
	removeNext  bool
}

func (tc *traceCleaner) Write(p []byte) (n int, err error) {
	tc.line += string(p)
OUTER:
	for {
		nextLineIndex := strings.Index(tc.line, "\n") + 1
		if nextLineIndex == 0 {
			return len(p), nil
		}
		line := tc.line[:nextLineIndex]
		tc.line = tc.line[nextLineIndex:]
		if tc.removeNext {
			n += len(line)
			tc.removeNext = false
			continue
		}
		for _, verboseRegExp := range verboseRegExps {
			if verboseRegExp.MatchString(line[:len(line)-1]) {
				n += len(line)
				tc.removeNext = true
				continue OUTER
			}
		}
		written, err := tc.destination.Write([]byte(line))
		n += written
		if err != nil {
			return n, err
		}
	}
}
