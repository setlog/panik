package panik

import (
	"fmt"
	"io"
	"strings"
)

type traceCleaner struct {
	destination io.Writer
	line        string
	removeNext  bool
}

func (tc *traceCleaner) Write(p []byte) (n int, err error) {
	tc.line += string(p)
OUTER:
	for nextLineIndex := strings.Index(tc.line, "\n") + 1; nextLineIndex > 0; nextLineIndex = strings.Index(tc.line, "\n") + 1 {
		line := tc.line[:nextLineIndex]
		tc.line = tc.line[nextLineIndex:]
		if tc.removeNext {
			n += len(line)
			tc.removeNext = false
			continue
		}
		for _, cleanFunc := range cleanFuncs {
			if strings.HasPrefix(line, fmt.Sprintf("%s(", cleanFunc)) {
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
	return len(p), nil
}
