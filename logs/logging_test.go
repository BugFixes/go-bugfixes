package logs_test

import (
	"errors"
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"testing"
)

func TestUnwrap(t *testing.T) {
	err := fmt.Errorf("%w error", errors.New("error"))
	wrappedError := logs.BugFixes{
		Err: err,
	}
	unwrappedError := wrappedError.UnwrapIt(err)

	if unwrappedError.Error() != "error" {
		t.Fatalf("Expected unwrapped error message to be 'error' instead got %v", unwrappedError.Error())
	}
}

func TestNewBugFixes(t *testing.T) {
	err := errors.New("error")
	bugFixesError := logs.NewBugFixes(err)

	if bugFixesError == nil {
		t.Fatal("Expected BugFixes error to be 'error'")
	}
}

func TestConvertLevelFromString(t *testing.T) {
	tests := []struct {
		input  string
		output int
	}{
		{"log", 1},
		{"debug", 1},
		{"info", 2},
		{"warn", 3},
		{"error", 4},
		{"crash", 5},
		{"panic", 5},
		{"fatal", 5},
		{"unknown", 9},
		{"10", 9},
		{"unrecognized", 9},
	}

	for _, test := range tests {
		converted := logs.ConvertLevelFromString(test.input)

		if converted != test.output {
			t.Fatalf("Expected '%v' to convert to %d, got %d", test.input, test.output, converted)
		}
	}
}
