package panik_test

import (
	"fmt"
	"testing"

	"github.com/setlog/panik"
)

func TestPanic(t *testing.T) {
	defer func() {
		r := recover()
		expectKnownCause(t, r)
		errMessage := r.(error).Error()
		if errMessage != "oof: 42 43" {
			t.Fatalf("error message was %v", errMessage)
		}
	}()
	panik.Panic("oof: ", 42, 43)
}

func TestPanicf(t *testing.T) {
	defer func() {
		r := recover()
		expectKnownCause(t, r)
		errMessage := r.(error).Error()
		if errMessage != "oof: 42" {
			t.Fatalf("error message was %v", errMessage)
		}
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

func expectUnknownCause(t *testing.T, r interface{}) {
	if r == nil {
		t.Fatal("r was nil")
	}
	if panik.HasKnownCause(r) {
		t.Fatal("r was a known cause")
	}
}

func TestOnError(t *testing.T) {
	var err error
	panik.OnError(err)
	err = fmt.Errorf("an error")
	defer func() {
		r := recover()
		expectKnownCause(t, r)
		errMessage := r.(error).Error()
		if errMessage != "an error" {
			t.Fatalf("error message was %v", errMessage)
		}
	}()
	panik.OnError(err)
}

func TestOnErrorf(t *testing.T) {
	var err error
	panik.OnErrorf(err, "oof")
	err = fmt.Errorf("an error")
	defer func() {
		r := recover()
		expectKnownCause(t, r)
		errMessage := r.(error).Error()
		if errMessage != "oof: 42: an error" {
			t.Fatalf("error message was %v", errMessage)
		}
	}()
	panik.OnErrorf(err, "oof: %d", 42)
}

func TestIfError(t *testing.T) {
	var err error
	panik.IfError(err, fmt.Errorf("foo"))
	err = fmt.Errorf("an error")
	defer func() {
		r := recover()
		expectKnownCause(t, r)
		errMessage := r.(error).Error()
		if errMessage != "foo" {
			t.Fatalf("error message was %v", errMessage)
		}
	}()
	panik.IfError(err, fmt.Errorf("foo"))
}

func TestIfErrorf(t *testing.T) {
	var err error
	panik.IfErrorf(err, "bla")
	err = fmt.Errorf("an error")
	defer func() {
		r := recover()
		expectKnownCause(t, r)
		errMessage := r.(error).Error()
		if errMessage != "foo: an error" {
			t.Fatalf("error message was %v", errMessage)
		}
	}()
	panik.IfErrorf(err, "foo: %w", panik.Cause{})
}

func TestWrap(t *testing.T) {
	defer func() {
		r := recover()
		expectUnknownCause(t, r)
		err, isError := r.(error)
		if !isError {
			t.Fatal("received non-error")
		}
		errMessage := err.Error()
		expected := "C42 43: panic"
		if errMessage != expected {
			t.Fatalf("Error() returned \"%s\". Expected \"%s\".", errMessage, expected)
		}
	}()
	defer panik.Wrap("C", 42, 43)
	panic("panic")
}

func TestWrapf(t *testing.T) {
	defer func() {
		r := recover()
		expectUnknownCause(t, r)
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
	defer panik.Wrapf("C%d", 42)
	panic("panic")
}

func TestToErrorCatchesKnownError(t *testing.T) {
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

func TestToErrorDoesNotCatchUnknownError(t *testing.T) {
	defer func() {
		r := recover()
		expectUnknownCause(t, r)
	}()
	dontCatchPanic()
}

func catchPanic() (retErr error) {
	defer panik.ToError(&retErr)
	defer panik.Wrapf("a: %d", 42)
	panik.Panicf("oof")
	return retErr
}

func dontCatchPanic() (retErr error) {
	defer panik.ToError(&retErr)
	defer panik.Wrapf("a: %d", 42)
	panic("oof")
}

func newCustomError(cause error, args ...interface{}) error {
	return fmt.Errorf("custom error %d: %w", args[0], cause)
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
	defer panik.Handle(func(r interface{}) {
		if r == nil {
			t.Fatal("handler was called with nil error")
		}
		handled = true
	})
	panic(catchPanic())
}

func TestHandlePanicsAgainOnUnknownValue(t *testing.T) {
	handled := false
	defer func() {
		if !handled {
			t.Fatalf("unknown panic was not handled")
		}
		r := recover()
		if r == nil {
			t.Fatal("unknown panic value was not thrown again")
		}
	}()
	defer panik.Handle(func(r interface{}) {
		handled = true
		if r != 42 {
			t.Fatal("r was not 42")
		}
	})
	panic(42)
}

func TestHandleConsumesKnownValue(t *testing.T) {
	handled := false
	defer func() {
		if !handled {
			t.Fatalf("known panic was not handled")
		}
		r := recover()
		if r != nil {
			t.Fatal("known panic value was thrown again")
		}
	}()
	defer panik.Handle(func(r interface{}) {
		handled = true
	})
	panik.Panic(42)
}
