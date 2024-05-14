package logs

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func Local(skipDepthOverride ...int) *BugFixes {
	if len(skipDepthOverride) > 0 && skipDepthOverride[0] != 0 {
		return &BugFixes{
			LocalOnly:         true,
			SkipDepthOverride: skipDepthOverride[0],
		}
	}

	return &BugFixes{
		LocalOnly: true,
	}
}

// <editor-fold desc="Error">
func (b *BugFixes) Error() string {
	return fmt.Sprintf("%s: %s", b.Bug, b.Err.Error())
}

func Error(inputs ...interface{}) error {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return Errorf(format, inputs...)
}
func Errorf(format string, inputs ...interface{}) error {
	b := &BugFixes{
		LocalOnly: false,
	}
	return b.Errorf(format, inputs...)
}
func (b *BugFixes) Errorf(format string, inputs ...interface{}) error {
	b.Level = "error"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.FormattedError = fmt.Errorf(format, inputs...)

	if !b.LocalOnly {
		b.Stack = debug.Stack()
		b.DoReporting()
	}

	return fmt.Errorf(format, inputs...)
}

// </editor-fold>

// Info <editor-fold desc="Info">
func Info(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return Infof(format, inputs...)
}
func (b *BugFixes) Info(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return b.Infof(format, inputs...)
}
func Infof(format string, inputs ...interface{}) string {
	b := &BugFixes{
		LocalOnly: false,
	}
	return b.Infof(format, inputs...)
}
func (b *BugFixes) Infof(format string, inputs ...interface{}) string {
	b.Level = "info"
	b.FormattedLog = fmt.Sprintf(format, inputs...)

	if !b.LocalOnly {
		b.Stack = debug.Stack()
		b.DoReporting()
	}

	return fmt.Sprintf("Info: %s", fmt.Sprintf(format, inputs...))
}

// </editor-fold>

// Debug <editor-fold desc="Debug">
func Debug(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return Debugf(format, inputs...)
}
func (b *BugFixes) Debug(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return b.Debugf(format, inputs...)
}
func Debugf(format string, inputs ...interface{}) string {
	b := &BugFixes{
		LocalOnly: false,
	}

	return b.Debugf(format, inputs...)
}
func (b *BugFixes) Debugf(format string, inputs ...interface{}) string {
	b.Level = "debug"
	b.FormattedLog = fmt.Sprintf(format, inputs...)

	if !b.LocalOnly {
		b.Stack = debug.Stack()
		b.DoReporting()
	}

	return fmt.Sprintf("Debug: %s", fmt.Sprintf(format, inputs...))
}

// </editor-fold>

// Log <editor-fold desc="Log">
func Log(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return Logf(format, inputs...)
}
func (b *BugFixes) Log(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return b.Logf(format, inputs...)
}
func Logf(format string, inputs ...interface{}) string {
	b := &BugFixes{
		LocalOnly: false,
	}
	return b.Logf(format, inputs...)
}
func (b *BugFixes) Logf(format string, inputs ...interface{}) string {
	b.Level = "log"
	b.FormattedLog = fmt.Sprintf(format, inputs...)

	if !b.LocalOnly {
		b.Stack = debug.Stack()
		b.DoReporting()
	}

	return fmt.Sprintf("Log: %s", fmt.Sprintf(format, inputs...))
}

// </editor-fold>

// Warn <editor-fold desc="Warn">
func Warn(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return Warnf(format, inputs...)
}
func (b *BugFixes) Warn(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return b.Warnf(format, inputs...)
}
func Warnf(format string, inputs ...interface{}) string {
	b := &BugFixes{
		LocalOnly: false,
	}
	return b.Warnf(format, inputs...)
}
func (b *BugFixes) Warnf(format string, inputs ...interface{}) string {
	b.Level = "warn"
	b.FormattedLog = fmt.Sprintf(format, inputs...)

	if !b.LocalOnly {
		b.Stack = debug.Stack()
		b.DoReporting()
	}

	return fmt.Sprintf("Warn: %s", fmt.Sprintf(format, inputs...))
}

// </editor-fold>

func Fatal(inputs ...interface{}) {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ")
	Fatalf(format, inputs...)
}
func (b *BugFixes) Fatal(inputs ...interface{}) {
	format := strings.Repeat("%v ", len(inputs))
	format = strings.TrimRight(format, ", ")
	b.Fatalf(format, inputs...)
}
func Fatalf(format string, inputs ...interface{}) {
	b := &BugFixes{
		LocalOnly: false,
	}
	b.Fatalf(format, inputs...)
}
func (b *BugFixes) Fatalf(format string, inputs ...interface{}) {
	b.Level = "fatal"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.Stack = debug.Stack()

	if !b.LocalOnly {
		b.DoReporting()
	}

	panic(b)
}
