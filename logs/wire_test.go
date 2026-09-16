package logs

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWirePayloadHasNoCredentials(t *testing.T) {
	b := &BugFixes{FormattedLog: "x", Level: "error", AgentID: "the-key", Secret: "the-secret"}
	out, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "the-secret") || strings.Contains(string(out), "the-key") {
		t.Fatalf("credentials leaked into payload: %s", out)
	}
}
