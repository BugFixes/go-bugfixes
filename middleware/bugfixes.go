package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	bugfixes "github.com/bugfixes/go-bugfixes"
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

// BugFixes will create a new middleware handler from a http.Handler.
func BugFixes() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
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
	return "bugfixes/middleware context value " + k.name
}

func SendToBugfixes(rvr interface{}) {
	cfg := bugfixes.GetDefaultConfig()
	if cfg.AgentKey == "" || cfg.AgentSecret == "" {
		return
	}

	s := &System{
		Config: &cfg,
	}
	s.SendToBugfixes(rvr)
}

func (s *System) SendToBugfixes(rvr interface{}) {
	stack := debug.Stack()
	go s.sendToBugfixes(rvr, stack)
}

func (s *System) sendToBugfixes(rvr interface{}, debugStack []byte) {
	cfg := s.config()
	p := prettyStack{}
	bug, err := p.bugParse(debugStack, rvr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bugfixes: failed to parse bug: %v\n", err)
		return
	}

	body, err := json.Marshal(bug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bugfixes: failed to marshall bug: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), bugfixes.DefaultTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "POST", cfg.BugEndpoint(), bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "bugfixes: failed to new request: %v\n", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", cfg.AgentKey)
	request.Header.Set("X-API-SECRET", cfg.AgentSecret)

	client := cfg.GetHTTPClient()
	resp, err := client.Do(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bugfixes: failed to send bug: %v\n", err)
		return
	}
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "bugfixes: failed to close body: %v\n", err)
		}
	}
}

func (s *System) config() bugfixes.Config {
	cfg := bugfixes.GetDefaultConfig()
	if s != nil && s.Config != nil {
		cfg = cfg.Merge(*s.Config)
	}
	if s != nil {
		cfg = cfg.Merge(bugfixes.Config{
			AgentKey:    s.AgentID,
			AgentSecret: s.Secret,
		})
	}

	return cfg
}

func ParseBugLine(bugLine string) (string, string, int, error) {
	i := strings.Index(bugLine, ":")
	if i < 0 {
		return "", "", 0, fmt.Errorf("failed to find ':' in bug line: %s", bugLine)
	}
	j := strings.Index(bugLine, " ")
	if j < 0 {
		return "", "", 0, fmt.Errorf("failed to find ' ' in bug line: %s", bugLine)
	}
	file := bugLine[:i]
	lne := bugLine[i+1 : j]
	line, err := strconv.Atoi(lne)
	if err != nil {
		return file, lne, 0, fmt.Errorf("failed to convert line number: %w", err)
	}

	return file, lne, line, nil
}
