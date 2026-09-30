package logs

import "testing"

func TestLevelCapturesStackOnlyForErrors(t *testing.T) {
	for _, level := range []string{"debug", "log", "info", "warn", "crash", "panic", "fatal"} {
		if levelCapturesStack(level) {
			t.Errorf("levelCapturesStack(%q) = true, want false", level)
		}
	}
	if !levelCapturesStack("error") {
		t.Fatal(`levelCapturesStack("error") = false, want true`)
	}
}
