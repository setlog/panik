package panik_test

import (
	"fmt"
	"testing"

	"github.com/setlog/panik"
)

func TestDescribe(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("received no panic")
		}
		err, isError := r.(error)
		if !isError {
			t.Fatal("received non-error")
		}
		errMessage := err.Error()
		expected := "C42: panic"
		if errMessage != expected {
			t.Fatalf("Error() returned \"%s\". Expected \"%s\".", errMessage, expected)
		}
	}()
	defer panik.Errorf("C%d: %w", 42, panik.Cause{})
	panic("panic")
}

func TestToError(t *testing.T) {
	err := catchPanic()
	if err == nil {
		t.Fatalf("err was nil")
	}
	message := err.Error()
	expectedMessage := "a: 42: oof"
	if message != expectedMessage {
		t.Fatalf("Message was \"%v\". Expected \"%v\".", message, expectedMessage)
	}
	if !panik.HasKnownCause(err) {
		t.Fatalf("err is not a known cause")
	}
}

func catchPanic() (retErr error) {
	defer panik.ToError(&retErr)
	defer panik.Errorf("a: %d: %w", 42, panik.Cause{})
	panik.Panicf("oof")
	return retErr
}

func newCustomError(cause error, args ...interface{}) error {
	return fmt.Errorf("custom error %d: %w", args[0], cause)
}

func TestOnError(t *testing.T) {
	var err error
	panik.OnError(err)
	err = fmt.Errorf("an error")
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("did not panic")
		}
	}()
	panik.OnError(err)
}

func TestPanicf(t *testing.T) {
	defer func() {
		expectKnownCause(t, recover())
	}()
	panik.Panicf("oof: %d", 42)
}

func expectKnownCause(t *testing.T, r interface{}) {
	if r == nil {
		t.Fatal("r was nil")
	}
	var err error
	var isError bool
	if err, isError = r.(error); !isError {
		t.Fatal("r was not an error")
	}
	if !panik.HasKnownCause(err) {
		t.Fatal("r was not a known cause")
	}
}

func TestHasKnownCause(t *testing.T) {
	if panik.HasKnownCause(nil) {
		t.Fatal("nil was a known cause")
	}
	err := catchPanic()
	if !panik.HasKnownCause(err) {
		t.Fatal("not a known cause")
	}
	err2 := fmt.Errorf("wrapped: %w", err)
	if !panik.HasKnownCause(err2) {
		t.Fatal("wrapped is not a known cause")
	}
}

func TestHandleReactsToKnownError(t *testing.T) {
	handled := false
	defer func() {
		if !handled {
			t.Error("handler was not called")
		}
		r := recover()
		if r != nil {
			t.Fatal("panic was not recovered")
		}
	}()
	defer panik.Handle(func(r error) {
		if r == nil {
			t.Fatal("handler was called with nil error")
		}
		handled = true
	})
	panic(catchPanic())
}

func TestHandleIgnoresValue(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("unknown panic was recovered")
		}
	}()
	defer panik.Handle(func(r error) {
		t.Fatalf("panik.Handle() reacted to unknown value with error %v", r)
	})
	panic(42)
}

func TestHandleIgnoresUnknownError(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("unknown panic was recovered")
		}
	}()
	defer panik.Handle(func(r error) {
		t.Fatalf("panik.Handle() reacted to unknown error %v", r)
	})
	panic(fmt.Errorf("oof"))
}
