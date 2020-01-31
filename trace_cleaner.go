package panik

import (
	"io"
	"regexp"
	"strings"
)

type traceCleaner struct {
	destination    io.Writer
	buffer         string
	removeNextLine bool
}

func (tc *traceCleaner) Write(p []byte) (n int, err error) {
	tc.buffer += string(p)
	for {
		nextLineIndex := strings.Index(tc.buffer, "\n") + 1
		if nextLineIndex == 0 {
			return len(p), nil
		}
		line := tc.buffer[:nextLineIndex]
		tc.buffer = tc.buffer[nextLineIndex:]
		if tc.removeNextLine {
			n += len(line)
			tc.removeNextLine = false
			continue
		}

		if isUnwantedLine(line) {
			n += len(line)
			tc.removeNextLine = true
		} else {
			written, err := tc.destination.Write([]byte(line))
			n += written
			if err != nil {
				return n, err
			}
		}
	}
}

var unwantedLineRegExps []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`^panic\(.*$`),
	regexp.MustCompile(`^runtime/debug.Stack\(.*$`),
	regexp.MustCompile(`^github.com/setlog/panik\..*\(.*$`),
}

func isUnwantedLine(line string) bool {
	for _, verboseRegExp := range unwantedLineRegExps {
		if verboseRegExp.MatchString(line[:len(line)-1]) {
			return true
		}
	}
	return false
}
