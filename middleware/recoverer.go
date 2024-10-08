package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns an HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
//
// Alternatively, look at https://github.com/pressly/lg middleware pkgs.
func Recoverer(next http.Handler) http.Handler {
	s := &System{}
	return s.Recoverer(next)
}

func (s *System) Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rvr := recover()
			if rvr != nil && rvr != http.ErrAbortHandler {

				logEntry := GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr)
				} else {
					PrintPrettyStack(rvr)
				}

				w.WriteHeader(http.StatusInternalServerError)

				go s.SendToBugfixes(rvr, http.Client{
					Timeout: time.Second * 10,
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func PrintPrettyStack(rvr interface{}) {
	debugStack := debug.Stack()
	s := prettyStack{}
	out, err := s.parse(debugStack, rvr)
	if err == nil {
		if _, errs := os.Stderr.Write(out); errs != nil {
			log.Fatal(errs)
		}
	} else {
		// print stdlib output as a fallback
		if _, errs := os.Stderr.Write(debugStack); errs != nil {
			log.Fatal(errs)
		}
	}
}

type prettyStack struct {
}

func (s prettyStack) parse(debugStack []byte, rvr interface{}) ([]byte, error) {
	var err error

	buf := &bytes.Buffer{}

	cW(buf, false, bRed, "\n")
	cW(buf, true, bCyan, " panic: ")
	cW(buf, true, bBlue, "%v", rvr)
	cW(buf, false, bWhite, "\n \n")

	// process debug stack info
	stack := strings.Split(string(debugStack), "\n")
	var lines []string

	// locate panic line, as we may have nested panics
	for i := len(stack) - 1; i > 0; i-- {
		lines = append(lines, stack[i])
		if strings.HasPrefix(stack[i], "panic(0x") {
			lines = lines[0 : len(lines)-2] // remove boilerplate
			break
		}
	}

	// reverse
	for i := len(lines)/2 - 1; i >= 0; i-- {
		opp := len(lines) - 1 - i
		lines[i], lines[opp] = lines[opp], lines[i]
	}

	// decorate
	for i, line := range lines {
		lines[i], err = s.decorateLine(line, true, i)
		if err != nil {
			return nil, err
		}
	}

	for _, l := range lines {
		if _, errs := fmt.Fprintf(buf, "%s", l); errs != nil {
			return nil, errs
		}
	}
	return buf.Bytes(), nil
}

func (s prettyStack) bugParse(debugStack []byte, rvr interface{}) (BugFixesSend, error) {
	bug := BugFixesSend{}

	var err error
	buf := &bytes.Buffer{}

	cW(buf, false, bRed, "\n")
	cW(buf, true, bCyan, " panic: ")
	cW(buf, true, bBlue, "%v", rvr)
	cW(buf, false, bWhite, "\n \n")

	// process debug stack info
	stack := strings.Split(string(debugStack), "\n")
	var lines []string

	bug.Level = "unknown"
	// locate panic line, as we may have nested panics
	for i := len(stack) - 1; i > 0; i-- {
		lines = append(lines, stack[i])
		if strings.HasPrefix(stack[i], "panic") {
			bug.Level = "panic"
			lines = lines[0 : len(lines)-2] // remove boilerplate
			break
		}
	}

	// reverse
	for i := len(lines)/2 - 1; i >= 0; i-- {
		opp := len(lines) - 1 - i
		lines[i], lines[opp] = lines[opp], lines[i]
	}

	bugLine := lines[1]
	i := strings.Index(bugLine, " ")
	bug.BugLine = bugLine[:i]

	file, line, lineNumber, err := ParseBugLine(lines[1])
	if err != nil {
		return bug, fmt.Errorf("failed to parse bug line: %w", err)
	}

	bug.Raw = flatten(lines, "\n")

	// decorate
	for i, line := range lines {
		lines[i], err = s.decorateLine(line, true, i)
		if err != nil {
			return bug, err
		}
	}

	bug.File = file
	bug.LineNumber = lineNumber
	bug.Line = line
	bug.Bug = flatten(lines, "")

	return bug, nil
}

func flatten(lines []string, seperator string) string {
	return strings.Join(lines[:], seperator)
}

func (s prettyStack) decorateLine(line string, useColor bool, num int) (string, error) {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "\t") || strings.Contains(line, ".go:") {
		return s.decorateSourceLine(line, useColor, num)
	} else if strings.HasSuffix(line, ")") {
		return s.decorateFuncCallLine(line, useColor, num)
	} else {
		if strings.HasPrefix(line, "\t") {
			return strings.Replace(line, "\t", "      ", 1), nil
		} else {
			return fmt.Sprintf("    %s\n", line), nil
		}
	}
}

func (s prettyStack) decorateFuncCallLine(line string, useColor bool, num int) (string, error) {
	idx := strings.LastIndex(line, "(")
	if idx < 0 {
		return "", errors.New("not a func call line")
	}

	buf := &bytes.Buffer{}
	pkg := line[0:idx]
	// addr := line[idx:]
	method := ""

	idx = strings.LastIndex(pkg, string(os.PathSeparator))
	if idx <= 0 {
		idx = strings.Index(pkg, ".")
		if idx <= 0 {
			method = pkg[0:]
			pkg = pkg[0:]
		} else {
			method = pkg[idx:]
			pkg = pkg[0:idx]
		}
	} else {
		method = pkg[idx+1:]
		pkg = pkg[0:idx+1]
		idx = strings.Index(method, ".")
		pkg += method[0:idx]
		method = method[idx:]
	}
	pkgColor := nYellow
	methodColor := bGreen

	if num == 0 {
		cW(buf, useColor, bRed, " -> ")
		pkgColor = bMagenta
		methodColor = bRed
	} else {
		cW(buf, useColor, bWhite, "    ")
	}
	cW(buf, useColor, pkgColor, "%s", pkg)
	cW(buf, useColor, methodColor, "%s\n", method)
	// cW(buf, useColor, nBlack, "%s", addr)
	return buf.String(), nil
}

func (s prettyStack) decorateSourceLine(line string, useColor bool, num int) (string, error) {
	idx := strings.LastIndex(line, ".go:")
	if idx < 0 {
		return "", errors.New("not a source line")
	}

	buf := &bytes.Buffer{}
	path := line[0:idx+3]
	lineno := line[idx+3:]

	idx = strings.LastIndex(path, string(os.PathSeparator))
	dir := path[0:idx+1]
	file := path[idx+1:]

	idx = strings.Index(lineno, " ")
	if idx > 0 {
		lineno = lineno[0:idx]
	}
	fileColor := bCyan
	lineColor := bGreen

	if num == 1 {
		cW(buf, useColor, bRed, " ->   ")
		fileColor = bRed
		lineColor = bMagenta
	} else {
		cW(buf, false, bWhite, "      ")
	}
	cW(buf, useColor, bWhite, "%s", dir)
	cW(buf, useColor, fileColor, "%s", file)
	cW(buf, useColor, lineColor, "%s", lineno)
	if num == 1 {
		cW(buf, false, bWhite, "\n")
	}
	cW(buf, false, bWhite, "\n")

	return buf.String(), nil
}
