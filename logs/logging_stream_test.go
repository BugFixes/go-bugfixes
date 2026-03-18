package logs

import (
	"io"
	"os"
	"strings"
	"testing"

	bugfixes "github.com/bugfixes/go-bugfixes"
)

func captureStandardStreams(t *testing.T, fn func()) (string, string) {
	t.Helper()

	origStdout := os.Stdout
	origStderr := os.Stderr

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stderr pipe: %v", err)
	}

	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	defer func() {
		os.Stdout = origStdout
		os.Stderr = origStderr
	}()

	fn()

	_ = stdoutWriter.Close()
	_ = stderrWriter.Close()

	stdout, err := io.ReadAll(stdoutReader)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	stderr, err := io.ReadAll(stderrReader)
	if err != nil {
		t.Fatalf("read stderr: %v", err)
	}

	return string(stdout), string(stderr)
}

func TestDoReportingLocalStreamsByLevel(t *testing.T) {
	tests := []struct {
		name         string
		level        string
		expectStdout bool
		expectStderr bool
	}{
		{name: "info uses stdout", level: INFO, expectStdout: true},
		{name: "warn uses stderr", level: WARN, expectStderr: true},
		{name: "error uses stderr", level: ERROR, expectStderr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logLine := "stream routing test"
			entry := &BugFixes{
				FormattedLog: logLine,
				Level:        test.level,
				Config: &bugfixes.Config{
					LocalOnly: true,
				},
			}

			stdout, stderr := captureStandardStreams(t, entry.DoReporting)

			if test.expectStdout && !strings.Contains(stdout, logLine) {
				t.Fatalf("expected stdout to contain %q, got %q", logLine, stdout)
			}
			if !test.expectStdout && stdout != "" {
				t.Fatalf("expected stdout to be empty, got %q", stdout)
			}
			if test.expectStderr && !strings.Contains(stderr, logLine) {
				t.Fatalf("expected stderr to contain %q, got %q", logLine, stderr)
			}
			if !test.expectStderr && stderr != "" {
				t.Fatalf("expected stderr to be empty, got %q", stderr)
			}
		})
	}
}
