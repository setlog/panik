package panik

import "testing"

func TestMakeErrorAddsFormattingDirective(t *testing.T) {
	err := makeError("%s no", 42, "oh")
	if err.Error() != "oh no: 42" {
		t.Fatalf("unexpected result: %v", err)
	}
}

func TestMakeErrorIsNotTrickedByPercentEscape(t *testing.T) {
	err := makeError("%s no: %%w", 42, "oh")
	if err.Error() != "oh no: %w: 42" {
		t.Fatalf("unexpected result: %v", err)
	}
}

func TestMakeErrorSeesExistingFormattingDirective(t *testing.T) {
	err := makeError("%s no: %w", 42, "oh", Cause{})
	if err.Error() != "oh no: 42" {
		t.Fatalf("unexpected result: %v", err)
	}
}

func TestHasErrorFormattingDirective(t *testing.T) {
	tests := []struct {
		format   string
		expected bool
	}{
		{"%w", true},
		{"%%w", false},
		{"%%%w", true},
		{"oof: %w", true},
		{"oof: %%w", false},
		{"oof: %%%w", true},
	}
	for i, test := range tests {
		result := hasErrorFormattingDirective.MatchString(test.format)
		if result != test.expected {
			t.Errorf("%d: (%v) yielded %v. Expected %v.", i, test.format, result, test.expected)
		}
	}
}
