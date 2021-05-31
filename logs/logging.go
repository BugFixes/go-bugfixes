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
)

type BugFixesLog struct {
	Log        string `json:"log"`
	Level      string `json:"level"`
	File       string `json:"file"`
	Line       string `json:"line"`
	LineNumber int    `json:"line_number"`
	Stack      []byte `json:"-"`
}

func (b BugFixesLog) DoReporting() {
	_, file, line, _ := runtime.Caller(2)
	b.File = file
	b.LineNumber = line
	b.Line = strconv.Itoa(line)

	keepLocal := os.Getenv("BUGFIXES_LOCAL_ONLY")
	if keepLocal == "" || keepLocal == "true" {
		b.makePretty()
		return
	}

	go func() {
		b.sendLog()
	}()
}

func (b BugFixesLog) sendLog() {
	agentKey := os.Getenv("BUGFIXES_AGENT_KEY")
	agentSecret := os.Getenv("BUGFIXES_AGENT_SECRET")

	bugServer := "https://api.bugfix.es/log"
	if bugServerEnv := os.Getenv("BUGFIXES_SERVER"); bugServerEnv != "" {
		bugServer = bugServerEnv
	}

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
		fmt.Printf("bugfixes sendLog marshal: %+v", err)
	}

	request, err := http.NewRequest("POST", bugServer, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("bugfixes sendLog newRequest: %+v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", agentKey)
	request.Header.Set("X-API-SECRET", agentSecret)

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Printf("bugfixes sendLog do: %+v", err)
	}
	if err := resp.Body.Close(); err != nil {
		fmt.Printf("bugfixes sendLog close: %+v", err)
	}
}

func (b BugFixesLog) makePretty() {
	if b.Stack != nil {
		PrintPrettyStack(b.Stack)
		return
	}

	out := &bytes.Buffer{}

	switch b.Level {
	case "warn":
		cW(out, true, bBlue, b.Level)
	case "info":
		cW(out, true, bYellow, b.Level)

	case "log":
		cW(out, true, bGreen, b.Level)
	case "debug":
		cW(out, true, bMagenta, b.Level)

	case "error":
		cW(out, true, bRed, b.Level)
	}

	fmt.Printf("%s >> %s:%d\n", out, b.File, b.LineNumber)
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
