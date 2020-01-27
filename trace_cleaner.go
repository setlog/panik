package panik

import "io"

type traceCleaner struct {
	destination io.Writer
}

func (tc *traceCleaner) Write(p []byte) (n int, err error) {
	// TODO: Actually clean up
	return tc.destination.Write(p)
}
