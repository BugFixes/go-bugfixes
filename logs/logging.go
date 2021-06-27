package logs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/go-logfmt/logfmt"
)

type BugFixesLog struct {
	FormattedLog string `json:"log"`
	Level        string `json:"level"`
	File         string `json:"file"`
	Line         string `json:"line"`
	LineNumber   int    `json:"line_number"`
	LogFmt       string `json:"log_fmt"`
	Stack        []byte `json:"stack"`

	FormattedError error `json:"-"`
	LocalOnly      bool  `json:"-"`
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

func GetLevelLog() int {
	return 1
}

func GetLevelInfo() int {
	return 2
}

func GetLevelError() int {
	return 3
}

func GetLevelCrash() int {
	return 4
}

func GetLevelUnknown() int {
	return 5
}

// ConvertLevelFromString
// nolint: gocyclo
func ConvertLevelFromString(s string) int {
	switch s {
	case LOG:
	case DEBUG:
		return GetLevelLog()

	case INFO:
	case WARN:
		return GetLevelInfo()

	case ERROR:
		return GetLevelError()

	case CRASH:
	case PANIC:
	case FATAL:
		return GetLevelCrash()

	case UNKNOWN:
		return GetLevelUnknown()

	default:
		lvl, err := strconv.Atoi(s)
		if err != nil {
			return GetLevelUnknown()
		}
		if lvl >= 5 {
			return GetLevelUnknown()
		}
		return lvl
	}

	return GetLevelUnknown()
}

func (b BugFixesLog) DoReporting() {
	skip := 2
	if b.LocalOnly {
		skip = 3
	}
	_, file, line, _ := runtime.Caller(skip)
	b.File = file
	b.LineNumber = line
	b.Line = strconv.Itoa(line)

	// Log Format
	b.logFormat()

	// Make it pretty
	b.makePretty()

	// Do we keep it local no matter what
	keepLocal := os.Getenv("BUGFIXES_LOCAL_ONLY")
	if keepLocal == "" || keepLocal == "true" || b.LocalOnly {
		return
	}

	// Log level
	reportLogLevel := ConvertLevelFromString(os.Getenv("BUGFIXES_LOG_LEVEL"))
	logLevel := ConvertLevelFromString(b.Level)
	if reportLogLevel > logLevel {
		return
	}

	b.sendLog()
}

func (b *BugFixesLog) logFormat() {
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

func (b BugFixesLog) sendLog() {
	agentKey := os.Getenv("BUGFIXES_AGENT_KEY")
	agentSecret := os.Getenv("BUGFIXES_AGENT_SECRET")

	bugServer := "https://api.bugfix.es"
	if bugServerEnv := os.Getenv("BUGFIXES_SERVER"); bugServerEnv != "" {
		bugServer = bugServerEnv
	}
	bugServer = fmt.Sprintf("%s/log", bugServer)

	if agentKey == "" || agentSecret == "" {
		fmt.Printf("cant send to server till you have created an agent and set the keys\n")
		if agentKey == "" {
			fmt.Printf("env: BUGFIXES_AGENT_KEY missing\n")
		}
		if agentSecret == "" {
			fmt.Printf("env: BUGFIXES_AGENT_SECRET missing\n")
		}
		return
	}

	body, err := json.Marshal(b)
	if err != nil {
		fmt.Printf("bugfixes sendLog marshal: %+v\n", err)
		return
	}

	request, err := http.NewRequest("POST", bugServer, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("bugfixes sendLog newRequest: %+v\n", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", agentKey)
	request.Header.Set("X-API-SECRET", agentSecret)

	client := http.Client{
		Timeout: 5 * time.Second,
	}
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

func (b BugFixesLog) makePretty() {
	out := &bytes.Buffer{}
	log := b.FormattedLog

	switch b.Level {
	case "warn":
		cW(out, true, bBlue, "Warning:")
	case "info":
		cW(out, true, bYellow, "Info:")

	case "log":
		cW(out, true, bGreen, "Log:")
	case "debug":
		cW(out, true, bMagenta, "Debug:")

	case "error":
		cW(out, true, bRed, "Error:")
		log = b.FormattedError.Error()

	default:
		cW(out, true, bWhite, fmt.Sprintf("%s:", b.Level))
	}

	fmt.Printf("%s %s:%d >> %s\n", out, b.File, b.LineNumber, log)

	if b.Stack != nil {
		extra := &bytes.Buffer{}
		cW(extra, true, bMagenta, "Stack:")
		fmt.Printf("%s", extra)
		PrintPrettyStack(b.Stack)
		return
	}
}

var (
	// Normal colors
	//nRed    = []byte{'\033', '[', '3', '1', 'm'}
	//nGreen  = []byte{'\033', '[', '3', '2', 'm'}
	nYellow = []byte{'\033', '[', '3', '3', 'm'}
	//nCyan   = []byte{'\033', '[', '3', '6', 'm'}
	// Bright colors
	bRed     = []byte{'\033', '[', '3', '1', ';', '1', 'm'}
	bGreen   = []byte{'\033', '[', '3', '2', ';', '1', 'm'}
	bYellow  = []byte{'\033', '[', '3', '3', ';', '1', 'm'}
	bBlue    = []byte{'\033', '[', '3', '4', ';', '1', 'm'}
	bMagenta = []byte{'\033', '[', '3', '5', ';', '1', 'm'}
	bCyan    = []byte{'\033', '[', '3', '6', ';', '1', 'm'}
	bWhite   = []byte{'\033', '[', '3', '7', ';', '1', 'm'}

	reset = []byte{'\033', '[', '0', 'm'}
)

var IsTTY bool

func init() {
	// This is sort of cheating: if stdout is a character device, we assume
	// that means it's a TTY. Unfortunately, there are many non-TTY
	// character devices, but fortunately stdout is rarely set to any of
	// them.
	//
	// We could solve this properly by pulling in a dependency on
	// code.google.com/p/go.crypto/ssh/terminal, for instance, but as a
	// heuristic for whether to print in color or in black-and-white, I'd
	// really rather not.
	fi, err := os.Stdout.Stat()
	if err == nil {
		m := os.ModeDevice | os.ModeCharDevice
		IsTTY = fi.Mode()&m == m
	}
}

// colorWrite
func cW(w io.Writer, useColor bool, color []byte, s string, args ...interface{}) {
	if IsTTY && useColor {
		_, _ = w.Write(color)
	}
	_, _ = fmt.Fprintf(w, s, args...)
	if IsTTY && useColor {
		_, _ = w.Write(reset)
	}
}
