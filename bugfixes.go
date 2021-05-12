package bugfixes

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
  "os"
  "runtime/debug"
  "time"
)

type BugFixes struct {
  Bug interface{} `json:"bug"`
  Raw interface{} `json:"raw"`
  BugLine string `json:"bug_line"`
  File string `json:"file"`
  Line int `json:"line"`
}

// New will create a new middleware handler from a http.Handler.
func New(h http.Handler) func(next http.Handler) http.Handler {
  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      h.ServeHTTP(w, r)
    })
  }
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
  name string
}

func (k *contextKey) String() string {
  return "chi/middleware context value " + k.name
}

func sendToBugfixes(rvr interface{}) {
  agentKey := os.Getenv("BUGFIXES_AGENT_KEY")
  agentSecret := os.Getenv("BUGFIXES_AGENT_SECRET")
  if agentKey == "" || agentSecret == "" {
    return
  }

  client := http.Client{
    Timeout: time.Duration(5 * time.Second),
  }

  debugStack := debug.Stack()
  s := prettyStack{}
  out := &bytes.Buffer{}
  bug, err := s.bugParse(debugStack, rvr)
  if err != nil {
    fmt.Fprintf(out, "bugfixes: failed to parse bug: %v", err)
    os.Stderr.Write(out.Bytes())
  }

  bugServer := "https://api.bugfix.es/bug"
  if bugServerEnv := os.Getenv("BUGFIXES_SERVER"); bugServerEnv != "" {
    bugServer = bugServerEnv
  }

  body, err := json.Marshal(bug)
  if err != nil {
    fmt.Fprintf(out, "bugfixes: failed to marshall bug: %v", err)
    os.Stderr.Write(out.Bytes())
  }
  request, err := http.NewRequest("POST", bugServer, bytes.NewBuffer(body))
  request.Header.Set("Content-type", "application/json")
  request.Header.Set("X-API-KEY", agentKey)
  request.Header.Set("X-API-SECRET", agentSecret)
  if err != nil {
    fmt.Fprintf(out, "bugfixes: failed to new request: %v", err)
    os.Stderr.Write(out.Bytes())
  }

  if _, err := client.Do(request); err != nil {
    fmt.Fprintf(out, "bugfixes: failed to send bug: %v", err)
    os.Stderr.Write(out.Bytes())
  }
}

func parseBugLine(bugLine string) (string, int, error) {
  return "", 0, nil
}