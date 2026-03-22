package logs

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"unicode/utf8"
)

type prettyStack struct {
}

func PrintPrettyStack(rvr interface{}) {
	debugStack, panicValue, showPanicValue := prettyStackInput(rvr)
	s := prettyStack{}
	out, err := s.parseWithOptions(debugStack, panicValue, showPanicValue)
	if err == nil {
		_, _ = os.Stderr.Write(out)
	} else {
		// print stdlib output as a fallback
		_, _ = os.Stderr.Write(debugStack)
	}
}

func (s prettyStack) parse(debugStack []byte, rvr interface{}) ([]byte, error) {
	return s.parseWithOptions(debugStack, rvr, true)
}

func (s prettyStack) parseWithOptions(debugStack []byte, rvr interface{}, showPanicValue bool) ([]byte, error) {
	var err error
	buf := &bytes.Buffer{}

	cW(buf, false, bRed, "\n")
	if showPanicValue {
		cW(buf, true, bCyan, " panic: ")
		cW(buf, true, bBlue, "%v", rvr)
		cW(buf, false, bWhite, "\n \n")
	}

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

func prettyStackInput(rvr interface{}) ([]byte, interface{}, bool) {
	switch v := rvr.(type) {
	case []byte:
		if looksLikeStackTrace(v) {
			return v, nil, false
		}
		return debug.Stack(), readableBytes(v), true
	case string:
		if looksLikeStackTrace([]byte(v)) {
			return []byte(v), nil, false
		}
		return debug.Stack(), v, true
	default:
		return debug.Stack(), rvr, true
	}
}

func looksLikeStackTrace(stack []byte) bool {
	return bytes.Contains(stack, []byte("goroutine ")) && bytes.Contains(stack, []byte(".go:"))
}

func readableBytes(value []byte) interface{} {
	if utf8.Valid(value) {
		return string(value)
	}

	return fmt.Sprintf("%x", value)
}

func (s prettyStack) decorateLine(line string, useColor bool, num int) (string, error) {
	line = strings.TrimSpace(line)
	switch {
	case strings.HasPrefix(line, "\t") || strings.Contains(line, ".go:"):
		return s.decorateSourceLine(line, useColor, num)
	case strings.HasSuffix(line, ")"):
		return s.decorateFuncCallLine(line, useColor, num)
	case strings.HasPrefix(line, "\t"):
		return strings.Replace(line, "\t", "      ", 1), nil
	default:
		return fmt.Sprintf("    %s\n", line), nil
	}
}

func (s prettyStack) decorateFuncCallLine(line string, useColor bool, num int) (string, error) {
	idx := strings.LastIndex(line, "(")
	if idx < 0 {
		return "", errors.New("not a func call line")
	}

	buf := &bytes.Buffer{}
	pkg := line[0:idx]
	method := ""

	idx = strings.LastIndex(pkg, string(os.PathSeparator))
	if idx < 0 {
		idx = strings.Index(pkg, ".")
		if idx < 0 {
			method = pkg
			pkg = ""
		} else {
			method = pkg[idx:]
			pkg = pkg[0:idx]
		}
	} else {
		method = pkg[idx+1:]
		pkg = pkg[0 : idx+1]
		idx = strings.Index(method, ".")
		if idx >= 0 {
			pkg += method[0:idx]
			method = method[idx:]
		}
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
	return buf.String(), nil
}

func (s prettyStack) decorateSourceLine(line string, useColor bool, num int) (string, error) {
	idx := strings.LastIndex(line, ".go:")
	if idx < 0 {
		return "", errors.New("not a source line")
	}

	buf := &bytes.Buffer{}
	path := line[0 : idx+3]
	lineno := line[idx+3:]

	idx = strings.LastIndex(path, string(os.PathSeparator))
	dir := path[0 : idx+1]
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
