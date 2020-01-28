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
