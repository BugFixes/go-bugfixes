package logs

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWirePayloadHasNoCredentials(t *testing.T) {
	b := &BugFixes{
		FormattedLog: "x",
		Level:        "error",
		AgentID:      "the-key",
		Secret:       "the-secret",
		CommitSHA:    strings.Repeat("a", 40),
		Release:      "api@2.4.0",
		Environment:  "development",
	}
	out, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "the-secret") || strings.Contains(string(out), "the-key") {
		t.Fatalf("credentials leaked into payload: %s", out)
	}
	for _, value := range []string{strings.Repeat("a", 40), "api@2.4.0", "development"} {
		if !strings.Contains(string(out), value) {
			t.Fatalf("deployment metadata %q missing from payload: %s", value, out)
		}
	}
}
