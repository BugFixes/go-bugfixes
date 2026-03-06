package logs

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func Local(skipDepthOverride ...int) *BugFixes {
	b := &BugFixes{LocalOnly: true}
	if len(skipDepthOverride) > 0 && skipDepthOverride[0] != 0 {
		b.SkipDepthOverride = skipDepthOverride[0]
	}
	return b
}

// variadicFormat builds a "%v, %v, ..." format string for variadic inputs.
func variadicFormat(inputs []interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	return strings.TrimRight(format, ", ")
}

// levelCapturesStack returns true for levels that include a stack trace.
func levelCapturesStack(level string) bool {
	switch level {
	case "error", "debug", "warn":
		return true
	default:
		return false
	}
}

// logAt is the shared implementation for string-returning log levels.
func (b *BugFixes) logAt(level, format string, inputs ...interface{}) string {
	if b == nil {
		return (&BugFixes{LocalOnly: false}).logAt(level, format, inputs...)
	}

	b.Level = level
	b.FormattedLog = fmt.Sprintf(format, inputs...)

	if !b.LocalOnly {
		if levelCapturesStack(level) {
			b.Stack = debug.Stack()
		}
		b.DoReporting()
	}

	display := strings.ToUpper(level[:1]) + level[1:]
	return fmt.Sprintf("%s: %s", display, b.FormattedLog)
}

// Error implements the error interface.
func (b *BugFixes) Error() string {
	if b.Err == nil {
		return b.Bug
	}
	return fmt.Sprintf("%s: %s", b.Bug, b.Err.Error())
}

// Error / Errorf — returns error, captures stack, sets FormattedError.
func Error(inputs ...interface{}) error {
	return Errorf(variadicFormat(inputs), inputs...)
}

func Errorf(format string, inputs ...interface{}) error {
	return (&BugFixes{LocalOnly: false}).Errorf(format, inputs...)
}

func (b *BugFixes) Errorf(format string, inputs ...interface{}) error {
	if b == nil {
		return Errorf(format, inputs...)
	}

	b.Level = "error"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.FormattedError = fmt.Errorf(format, inputs...)

	if !b.LocalOnly {
		b.Stack = debug.Stack()
		b.DoReporting()
	}

	return b.FormattedError
}

// Info / Infof
func Info(inputs ...interface{}) string  { return Infof(variadicFormat(inputs), inputs...) }
func (b *BugFixes) Info(inputs ...interface{}) string { return b.Infof(variadicFormat(inputs), inputs...) }
func Infof(format string, inputs ...interface{}) string { return (&BugFixes{LocalOnly: false}).Infof(format, inputs...) }
func (b *BugFixes) Infof(format string, inputs ...interface{}) string { return b.logAt("info", format, inputs...) }

// Debug / Debugf
func Debug(inputs ...interface{}) string  { return Debugf(variadicFormat(inputs), inputs...) }
func (b *BugFixes) Debug(inputs ...interface{}) string { return b.Debugf(variadicFormat(inputs), inputs...) }
func Debugf(format string, inputs ...interface{}) string { return (&BugFixes{LocalOnly: false}).Debugf(format, inputs...) }
func (b *BugFixes) Debugf(format string, inputs ...interface{}) string { return b.logAt("debug", format, inputs...) }

// Log / Logf
func Log(inputs ...interface{}) string  { return Logf(variadicFormat(inputs), inputs...) }
func (b *BugFixes) Log(inputs ...interface{}) string { return b.Logf(variadicFormat(inputs), inputs...) }
func Logf(format string, inputs ...interface{}) string { return (&BugFixes{LocalOnly: false}).Logf(format, inputs...) }
func (b *BugFixes) Logf(format string, inputs ...interface{}) string { return b.logAt("log", format, inputs...) }

// Warn / Warnf
func Warn(inputs ...interface{}) string  { return Warnf(variadicFormat(inputs), inputs...) }
func (b *BugFixes) Warn(inputs ...interface{}) string { return b.Warnf(variadicFormat(inputs), inputs...) }
func Warnf(format string, inputs ...interface{}) string { return (&BugFixes{LocalOnly: false}).Warnf(format, inputs...) }
func (b *BugFixes) Warnf(format string, inputs ...interface{}) string { return b.logAt("warn", format, inputs...) }

// Fatal / Fatalf — always captures stack, panics.
func Fatal(inputs ...interface{})            { Fatalf(variadicFormat(inputs), inputs...) }
func (b *BugFixes) Fatal(inputs ...interface{}) { b.Fatalf(variadicFormat(inputs), inputs...) }
func Fatalf(format string, inputs ...interface{}) { (&BugFixes{LocalOnly: false}).Fatalf(format, inputs...) }

func (b *BugFixes) Fatalf(format string, inputs ...interface{}) {
	if b == nil {
		Fatalf(format, inputs...)
		return
	}

	b.Level = "fatal"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.Stack = debug.Stack()

	if !b.LocalOnly {
		b.DoReporting()
	}

	panic(b)
}
