package panik

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestTraceCleaner(t *testing.T) {
	trace := []byte(`goroutine 1 [running]:
github.com/setlog/panik.Described(0x4ed443, 0x16, 0xc000049ee0, 0x1, 0x1)
	/home/developer/github/panik/panik.go:42 +0xda
panic(0x4ccbe0, 0xc00000c180)
	/usr/local/go/src/runtime/panic.go:679 +0x1b2
github.com/setlog/panik.Described(0x4ec782, 0x12, 0x0, 0x0, 0x0)
	/home/developer/github/panik/panik.go:42 +0xda
panic(0x4c3ac0, 0x506250)
	/usr/local/go/src/runtime/panic.go:679 +0x1b2
main.f(0x54, 0x0)
	/home/developer/github/panik/ex/main.go:38 +0xa3
main.getSomething(0x2a, 0x0, 0x0, 0x0, 0x0)
	/home/developer/github/panik/ex/main.go:16 +0xe9
main.main()
	/home/developer/github/panik/ex/main.go:9 +0x2e
`)
	var previousResult []byte = nil
	for _, bytesPerCall := range []int{-1, 1, 2, 3, 4, 5, 6, 7, 60, 61, 62, 63, 120, 121, 122, 123} {
		result := runTraceCleaner(t, trace, bytesPerCall)
		if previousResult != nil && !bytes.Equal(previousResult, result) {
			t.Fatalf("Result did not match result of previous run.")
		}
		if len(result) < 100 {
			t.Fatalf("Result is oddly short (%d):\n%s", len(result), string(result))
		}
		previousResult = result
	}
}

func TestTraceCleanerNilWrite(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	traceCleaner := &traceCleaner{destination: buf}
	write(t, traceCleaner, nil)
}

func TestTraceCleanerEmptyWrite(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	traceCleaner := &traceCleaner{destination: buf}
	write(t, traceCleaner, []byte{})
}

func runTraceCleaner(t *testing.T, trace []byte, bytesPerCall int) []byte {
	originalLineCount := len(strings.Split(string(trace), "\n"))
	buf := bytes.NewBuffer(nil)
	traceCleaner := &traceCleaner{destination: buf}
	if bytesPerCall <= 0 {
		write(t, traceCleaner, trace)
	} else {
		l := len(trace)
		for i := 0; i < l; i += bytesPerCall {
			limit := i + bytesPerCall
			if limit > l {
				limit = l
			}
			write(t, traceCleaner, trace[i:limit])
		}
	}
	cleanTrace := buf.String()
	lines := strings.Split(cleanTrace, "\n")
	actualLineCount := len(lines)
	expectedLineCount := 8

	// Expected:
	// goroutine 1 [running]:
	//  main.f(0x54, 0x0)
	//  	/home/developer/github/panik/ex/main.go:38 +0xa3
	//  main.getSomething(0x2a, 0x0, 0x0, 0x0, 0x0)
	//  	/home/developer/github/panik/ex/main.go:16 +0xe9
	//  main.main()
	//  	/home/developer/github/panik/ex/main.go:9 +0x2e

	if actualLineCount != expectedLineCount {
		t.Fatalf("For %d bytes per call: cleaned up trace has %d lines:\n%s\nExpected it to have %d lines. Original line count was %d.",
			bytesPerCall, actualLineCount, buf.String(), expectedLineCount, originalLineCount)
	}
	if lines[7] != "" {
		t.Fatalf("Last line was not an empty line.")
	}

	return []byte(cleanTrace)
}

func write(t *testing.T, w io.Writer, p []byte) {
	n, err := w.Write(p)
	if n != len(p) {
		t.Fatalf("Write() wrote %d bytes. Expected %d", n, len(p))
	}
	if err != nil {
		t.Fatalf("Write() reported error: %v. Expected nil.", err)
	}
}
