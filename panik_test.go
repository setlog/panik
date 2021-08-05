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
		if errMessage != "oof" {
			t.Fatalf("error message was %v", errMessage)
		}
	}()
	panik.Panic("oof")
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
	if !panik.Caused(err) {
		t.Fatal("r was not a known cause")
	}
}

func expectUnknownCause(t *testing.T, r interface{}) {
	if r == nil {
		t.Fatal("r was nil")
	}
	if panik.Caused(r) {
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

func TestOnErrore(t *testing.T) {
	var err error
	panik.OnErrore(err, fmt.Errorf("foo"))
	err = fmt.Errorf("an error")
	defer func() {
		r := recover()
		expectKnownCause(t, r)
		errMessage := r.(error).Error()
		if errMessage != "foo: an error" {
			t.Fatalf("error message was %v", errMessage)
		}
	}()
	panik.OnErrore(err, fmt.Errorf("foo"))
}

func TestOnErrorfw(t *testing.T) {
	var err error
	panik.OnErrorfw(err, "oof")
	err = fmt.Errorf("an error")
	defer func() {
		r := recover()
		expectKnownCause(t, r)
		errMessage := r.(error).Error()
		if errMessage != "oof: 42: an error" {
			t.Fatalf("error message was %v", errMessage)
		}
	}()
	panik.OnErrorfw(err, "oof: %d", 42)
}

func TestOnErrorfv(t *testing.T) {
	var err error
	panik.OnErrorfv(err, "bla")
	err = fmt.Errorf("an error")
	defer func() {
		r := recover()
		expectKnownCause(t, r)
		errMessage := r.(error).Error()
		if errMessage != "foo: 42: an error" {
			t.Fatalf("error message was %v", errMessage)
		}
	}()
	panik.OnErrorfv(err, "foo: %d", 42)
}

func TestOnErrorfvRetainsKnownCause(t *testing.T) {
	defer func() {
		r := recover()
		if !panik.Caused(r) {
			t.Fatalf("r not marked as caused by panik")
		}
	}()
	panik.OnErrorfv(catchPanic(), "onerror")
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
	if panik.Caused(err) {
		t.Fatalf("err is still wrapped after deescalation with ToError()")
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

func TestToErrorRetainsErrorIdentity(t *testing.T) {
	var retErr error = nil
	panikErr := fmt.Errorf("foo")
	defer func() {
		if retErr != panikErr {
			t.Fatal("retErr != panikErr")
		}
	}()
	defer panik.ToError(&retErr)
	panik.OnError(panikErr)
}

func TestHasKnownCause(t *testing.T) {
	if panik.Caused(nil) {
		t.Fatal("nil was a known cause")
	}
	err := catchPanic()
	if panik.Caused(err) {
		t.Fatal("err is still wrapped after deescalation with ToError()")
	}
	err = fmt.Errorf("foo")
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("recovered nil")
		}
		if r == err {
			t.Fatal("r == err")
		}
		if !panik.Caused(r) {
			t.Fatal("directly recovered panik error not caused by panik")
		}
		err2 := fmt.Errorf("wrapped: %w", err)
		if panik.Caused(err2) {
			t.Fatal("wrapped is still a known cause")
		}
	}()
	panik.OnError(err)
}
