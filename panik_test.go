package panik_test

import (
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
	defer panik.Described("A: %w", panik.Cause{})
	defer panik.Describe("B: %w", panik.Cause{})
	defer panik.Describe("C%d: %w", 42, panik.Cause{})
	panic("panic")
}

func TestMultipleDescribe(t *testing.T) {
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
		expected := "C42: F42: panic"
		if errMessage != expected {
			t.Fatalf("Error() returned \"%s\". Expected \"%s\".", errMessage, expected)
		}
	}()
	defer panik.Described("A: %w", panik.Cause{})
	defer panik.Describe("B: %w", panik.Cause{})
	defer panik.Describe("C%d: %w", 42, panik.Cause{})
	panFunc()
}

func panFunc() {
	defer panik.Described("D: %w", panik.Cause{})
	defer panik.Describe("E: %w", panik.Cause{})
	defer panik.Describe("F%d: %w", 42, panik.Cause{})
	panic("panic")
}
