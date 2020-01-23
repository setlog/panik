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
	// defer panik.ConsolidateAsIs()
	defer panik.Describe("A")
	defer panik.Describe("B")
	defer panik.Describe("C%d", 42)
	panic("panic")
}
