package middleware_test

import (
	"fmt"
	"github.com/bugfixes/go-bugfixes/middleware"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestBugfixes(t *testing.T) {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, errs := fmt.Fprint(w, "This is a test."); errs != nil {
			t.Fatalf("Could not write to response writer: %v", errs)
		}
	})
	handler := middleware.BugFixes(handlerFunc)
	testHandler := handler(handlerFunc)
	server := httptest.NewServer(handler(testHandler))
	resp, err := http.Get(server.URL)

	if err != nil {
		t.Fatalf("Could not send GET request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code 200, got %v", resp.StatusCode)
	}
}

func TestParseBugLine(t *testing.T) {
	file, lne, line, err := middleware.ParseBugLine("Example.go:53 def")

	expectedFile := "Example.go"
	expectedLne := "53"
	expectedLine := 53

	if err != nil {
		t.Fatalf("An error occurred: %v", err)
	}
	if file != expectedFile {
		t.Errorf("Expected %v, got %v", expectedFile, file)
	}
	if lne != expectedLne {
		t.Errorf("Expected %v, got %v", expectedLne, lne)
	}
	if line != expectedLine {
		t.Errorf("Expected %v, got %v", expectedLine, line)
	}
}

func TestSendToBugfixes(t *testing.T) {
	if err := os.Setenv("BUGFIXES_AGENT_KEY", "test_key"); err != nil {
		t.Fatalf("Could not set environment variable: %v", err)
	}
	if err := os.Setenv("BUGFIXES_AGENT_SECRET", "test_secret"); err != nil {
		t.Fatalf("Could not set environment variable: %v", err)
	}

	// Additional required configuration for the "sendToBugfixes" function might go here
	middleware.SendToBugfixes(nil)
}
