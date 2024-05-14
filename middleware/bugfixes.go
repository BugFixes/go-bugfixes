package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type BugFixesSend struct {
	Bug        interface{} `json:"bug"`
	Raw        interface{} `json:"raw"`
	BugLine    string      `json:"bug_line"`
	File       string      `json:"file"`
	Line       string      `json:"line"`
	LineNumber int         `json:"line_number"`
	Level      string      `json:"level"`
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

func SendToBugfixes(rvr interface{}) {
  agentKey := os.Getenv("BUGFIXES_AGENT_KEY")
  agentSecret := os.Getenv("BUGFIXES_AGENT_SECRET")
  if agentKey == "" || agentSecret == "" {
    return
  }

  s := &System{
    AgentID: agentKey,
    Secret: agentSecret,
  }
  s.SendToBugfixes(rvr)
}

func (s *System) SendToBugfixes(rvr interface{}) {
  client := http.Client{
		Timeout: 5 * time.Second,
	}

	debugStack := debug.Stack()
	p := prettyStack{}
	out := &bytes.Buffer{}
	bug, err := p.bugParse(debugStack, rvr)
	if err != nil {
		if _, errs := fmt.Fprintf(out, "bugfixes: failed to parse bug: %v", err); errs != nil {
			log.Fatal(errs)
		}
		if _, errs := os.Stderr.Write(out.Bytes()); errs != nil {
			log.Fatal(errs)
		}
		return
	}

	bugServer := "https://api.bugfix.es"
	if bugServerEnv := os.Getenv("BUGFIXES_SERVER"); bugServerEnv != "" {
		bugServer = bugServerEnv
	}
	bugServer = fmt.Sprintf("%s/bug", bugServer)

	body, err := json.Marshal(bug)
	if err != nil {
		if _, errs := fmt.Fprintf(out, "bugfixes: failed to marshall bug: %v", err); errs != nil {
			log.Fatal(errs)
		}
		if _, errs := os.Stderr.Write(out.Bytes()); errs != nil {
			log.Fatal(errs)
		}
		return
	}
	request, err := http.NewRequest("POST", bugServer, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", s.AgentID)
	request.Header.Set("X-API-SECRET", s.Secret)
	if err != nil {
		if _, errs := fmt.Fprintf(out, "bugfixes: failed to new request: %v", err); errs != nil {
			log.Fatal(errs)
		}
		if _, err := os.Stderr.Write(out.Bytes()); err != nil {
			log.Fatal(err)
		}
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		if _, errs := fmt.Fprintf(out, "bugfixes: failed to send bug: %v", err); errs != nil {
			log.Fatal(errs)
		}
		if _, errs := os.Stderr.Write(out.Bytes()); errs != nil {
			log.Fatal(errs)
		}
		return
	}
	if err := resp.Body.Close(); err != nil {
		if _, errs := fmt.Fprintf(out, "bugfixes: failed to close body: %v", err); errs != nil {
			log.Fatal(errs)
		}
		if _, errs := os.Stderr.Write(out.Bytes()); errs != nil {
			log.Fatal(errs)
		}
		return
	}
}

func ParseBugLine(bugLine string) (string, string, int, error) {
	i := strings.Index(bugLine, ":")
	j := strings.Index(bugLine, " ")
	file := bugLine[:i]
	lne := bugLine[i+1 : j]
	line, err := strconv.Atoi(lne)
	if err != nil {
		return file, lne, 0, fmt.Errorf("failed to convert line number: %w", err)
	}

	return file, lne, line, nil
}
