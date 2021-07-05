package logs

import (
	"fmt"
	"runtime/debug"
)

func Local() BugFixes {
	return BugFixes{
		LocalOnly: true,
	}
}

func (b BugFixes) Error() string {
	return fmt.Sprintf("%s: %s", b.Bug, b.Err.Error())
}

func Error(inputs ...interface{}) error {
	return Errorf("%w", inputs...)
}
func (b BugFixes) Errorf(format string, inputs ...interface{}) error {
	b.Level = "error"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.FormattedError = fmt.Errorf(format, inputs...)
	b.DoReporting()

	return fmt.Errorf(format, inputs...)
}
func Errorf(format string, inputs ...interface{}) error {
	b := BugFixes{
		Level:          "error",
		FormattedLog:   fmt.Sprintf(format, inputs...),
		FormattedError: fmt.Errorf(format, inputs...),
	}
	b.DoReporting()

	return fmt.Errorf(format, inputs...)
}

func (b BugFixes) Info(inputs ...interface{}) {
	b.Infof("%v", inputs...)
}
func Info(inputs ...interface{}) {
	Infof("%v", inputs...)
}
func (b BugFixes) Infof(format string, inputs ...interface{}) {
	b.Level = "info"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.DoReporting()
}
func Infof(format string, inputs ...interface{}) {
	b := BugFixes{
		Level:        "info",
		FormattedLog: fmt.Sprintf(format, inputs...),
	}
	b.DoReporting()
}

func (b BugFixes) Debug(inputs ...interface{}) {
	b.Debugf("%v", inputs...)
}
func Debug(inputs ...interface{}) {
	Debugf("%v", inputs...)
}
func (b BugFixes) Debugf(format string, inputs ...interface{}) {
	b.Level = "debug"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.Stack = debug.Stack()
	b.DoReporting()
}
func Debugf(format string, inputs ...interface{}) {
	b := BugFixes{
		Level:        "debug",
		FormattedLog: fmt.Sprintf(format, inputs...),
		Stack:        debug.Stack(),
	}
	b.DoReporting()
}

func (b BugFixes) Log(inputs ...interface{}) {
	b.Logf("%v", inputs...)
}
func Log(inputs ...interface{}) {
	Logf("%v", inputs...)
}
func (b BugFixes) Logf(format string, inputs ...interface{}) {
	b.Level = "log"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.DoReporting()
}
func Logf(format string, inputs ...interface{}) {
	b := BugFixes{
		Level:        "log",
		FormattedLog: fmt.Sprintf(format, inputs...),
	}
	b.DoReporting()
}

func (b BugFixes) Warn(inputs ...interface{}) {
	b.Warnf("%v", inputs...)
}
func Warn(inputs ...interface{}) {
	Warnf("%v", inputs...)
}
func (b BugFixes) Warnf(format string, inputs ...interface{}) {
	b.Level = "warn"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.DoReporting()
}
func Warnf(format string, inputs ...interface{}) {
	b := BugFixes{
		Level:        "warn",
		FormattedLog: fmt.Sprintf(format, inputs...),
	}
	b.DoReporting()
}
