package logs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	bugfixes "github.com/bugfixes/go-bugfixes"
	"github.com/bugfixes/go-bugfixes/internal/term"
	"github.com/go-logfmt/logfmt"
)

const logsPackagePrefix = "github.com/bugfixes/go-bugfixes/logs."

type BugFixes struct {
	FormattedLog string `json:"log"`
	Level        string `json:"level"`
	File         string `json:"file"`
	Line         string `json:"line"`
	LineNumber   int    `json:"line_number"`
	LogFmt       string `json:"log_fmt"`
	Stack        []byte `json:"stack"`

	FormattedError error `json:"-"`
	LocalOnly      bool  `json:"-"`

	Bug               string
	Err               error
	SkipDepthOverride int

	// Creds
	AgentID string
	Secret  string

	Config *bugfixes.Config `json:"-"`
}

func NewBugFixes(err error) error {
	if err == nil {
		return nil
	}

	return &BugFixes{
		Err: err,
	}
}

func (b *BugFixes) Setup(id, secret string) {
	b.AgentID = id
	b.Secret = secret
}

func (b *BugFixes) SetConfig(cfg bugfixes.Config) {
	b.Config = &cfg
}

const (
	LOG   = "log"
	DEBUG = "debug"

	INFO = "info"
	WARN = "warn"

	ERROR = "error"

	CRASH = "crash"
	PANIC = "panic"
	FATAL = "fatal"

	UNKNOWN = "unknown"
)

const (
	LevelDebug   = 1
	LevelLog     = 2
	LevelInfo    = 3
	LevelWarn    = 4
	LevelError   = 5
	LevelCrash   = 6
	LevelUnknown = 9
)

// ConvertLevelFromString converts a level name to its numeric value.
func ConvertLevelFromString(s string) int {
	switch s {
	case LOG:
		return LevelLog
	case DEBUG:
		return LevelDebug
	case INFO:
		return LevelInfo
	case WARN:
		return LevelWarn
	case ERROR:
		return LevelError
	case CRASH, PANIC, FATAL:
		return LevelCrash
	case UNKNOWN:
		return LevelUnknown
	default:
		lvl, err := strconv.Atoi(s)
		if err != nil {
			return LevelUnknown
		}
		if lvl >= LevelUnknown {
			return LevelUnknown
		}
		return lvl
	}
}

func (b *BugFixes) UnwrapIt(e error) error {
	u, ok := e.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}

	return u.Unwrap()
}

// findCaller walks the call stack and returns the first frame outside the logs package.
func (b *BugFixes) findCaller() {
	var pcs [25]uintptr
	n := runtime.Callers(1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		if !strings.HasPrefix(frame.Function, logsPackagePrefix) {
			b.File = frame.File
			b.LineNumber = frame.Line
			b.Line = strconv.Itoa(frame.Line)
			return
		}
		if !more {
			break
		}
	}
}

func (b *BugFixes) DoReporting() {
	cfg := b.config()

	b.findCaller()

	// Log Format
	b.logFormat()

	// Make it pretty
	b.makePretty()

	if cfg.LocalOnly {
		return
	}

	// Log level
	reportLogLevel := ConvertLevelFromString(cfg.LogLevel)
	logLevel := ConvertLevelFromString(b.Level)
	if reportLogLevel > logLevel {
		return
	}

	body, err := json.Marshal(b)
	if err != nil {
		fmt.Printf("bugfixes sendLog marshal: %+v\n", err)
		return
	}
	go b.sendLogBody(cfg, body)
}

func (b *BugFixes) logFormat() {
	out := bytes.Buffer{}
	lf := logfmt.NewEncoder(&out)

	if err := lf.EncodeKeyval("path", b.File); err != nil {
		fmt.Printf("logfmt path: %v", err)
	}
	if err := lf.EncodeKeyval("level", b.Level); err != nil {
		fmt.Printf("logfmt level: %v", err)
	}
	if err := lf.EncodeKeyval("msg", b.FormattedLog); err != nil {
		fmt.Printf("logfmt msg: %v", err)
	}
	if err := lf.EncodeKeyval("time", time.Now()); err != nil {
		fmt.Printf("logfmt time: %v", err)
	}
	if err := lf.EncodeKeyval("line", b.Line); err != nil {
		fmt.Printf("logfmt line: %v", err)
	}

	if err := lf.EndRecord(); err != nil {
		fmt.Printf("logfmt endrecord: %v", err)
	}

	b.LogFmt = out.String()
}

func (b *BugFixes) sendLogBody(cfg bugfixes.Config, body []byte) {
	if cfg.AgentKey == "" || cfg.AgentSecret == "" {
		fmt.Printf("cant send to server till you have created an agent and set the keys\n")
		if cfg.AgentKey == "" {
			fmt.Printf("env: BUGFIXES_AGENT_KEY missing\n")
		}
		if cfg.AgentSecret == "" {
			fmt.Printf("env: BUGFIXES_AGENT_SECRET missing\n")
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), bugfixes.DefaultTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "POST", cfg.LogEndpoint(), bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("bugfixes sendLog newRequest: %+v\n", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", cfg.AgentKey)
	request.Header.Set("X-API-SECRET", cfg.AgentSecret)

	client := cfg.GetHTTPClient()
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("bugfixes sendLog do: %+v\n", err)
		return
	}
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("bugfixes sendLog close: %+v\n", err)
			return
		}
	}
}

func (b *BugFixes) makePretty() {
	out := &bytes.Buffer{}
	log := b.FormattedLog
	cfg := b.config()

	switch b.Level {
	case "warn":
		cW(out, true, bYellow, "Warning:")
	case "info":
		cW(out, true, bCyan, "Info:")

	case "log":
		cW(out, true, bGreen, "Log:")
	case "debug":
		cW(out, true, bMagenta, "Debug:")

	case "error":
		cW(out, true, bRed, "Error:")
		if b.FormattedError != nil {
			log = b.FormattedError.Error()
		}

	default:
		cW(out, true, bWhite, "%s:", b.Level)
	}

	// print to stdout if the level is high enough
	reportLogLevel := ConvertLevelFromString(cfg.LogLevel)
	logLevel := ConvertLevelFromString(b.Level)
	if logLevel >= reportLogLevel || reportLogLevel == LevelUnknown || cfg.LocalOnly {
		fmt.Printf("%s %s >> %s:%d >> %s\n", out, time.Now().Format("2006-01-02 15:04:05"), b.File, b.LineNumber, log)
	}

	if b.Stack != nil {
		extra := &bytes.Buffer{}
		cW(extra, true, bMagenta, "Stack:")
		fmt.Printf("%s", extra)
		PrintPrettyStack(b.Stack)
		return
	}
}

func (b *BugFixes) config() bugfixes.Config {
	cfg := bugfixes.GetDefaultConfig()
	if b != nil && b.Config != nil {
		cfg = cfg.Merge(*b.Config)
	}
	if b != nil {
		cfg = cfg.Merge(bugfixes.Config{
			AgentKey:    b.AgentID,
			AgentSecret: b.Secret,
			LocalOnly:   b.LocalOnly,
		})
	}

	return cfg
}

var (
	nYellow  = term.NYellow
	bRed     = term.BRed
	bGreen   = term.BGreen
	bYellow  = term.BYellow
	bBlue    = term.BBlue
	bMagenta = term.BMagenta
	bCyan    = term.BCyan
	bWhite   = term.BWhite
)

// IsTTY reports whether stdout appears to be a terminal.
var IsTTY bool

func init() {
	IsTTY = term.IsTTY
}

// cW writes a color-formatted string to w.
func cW(w io.Writer, useColor bool, color []byte, s string, args ...interface{}) {
	term.CW(w, useColor, color, s, args...)
}
