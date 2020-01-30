package panik

import "testing"

func TestMakeErrorSeesExistingFormattingDirective(t *testing.T) {
	err := makeError("%s no: %w", makeCause(42), "oh", Cause{})
	if err.Error() != "oh no: 42" {
		t.Fatalf("unexpected result: %v", err)
	}
}
