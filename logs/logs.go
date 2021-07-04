package logs

import (
	"fmt"
	"runtime/debug"
)

func Local() BugFixesLog {
	return BugFixesLog{
		LocalOnly: true,
	}
}

func (b BugFixesLog) Unwrap(e error) error {
  u, ok := e.(interface {
    Unwrap() error
  })
  if !ok {
    return nil
  }

  return u.Unwrap()
}

func (b BugFixesLog) Error(inputs ...interface{}) error {
	return b.Errorf("%w", inputs...)
}
func Error(inputs ...interface{}) error {
	return Errorf("%w", inputs...)
}
func (b BugFixesLog) Errorf(format string, inputs ...interface{}) error {
	b.Level = "error"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.FormattedError = fmt.Errorf(format, inputs...)
	b.DoReporting()

	return fmt.Errorf(format, inputs...)
}
func Errorf(format string, inputs ...interface{}) error {
	b := BugFixesLog{
		Level:          "error",
		FormattedLog:   fmt.Sprintf(format, inputs...),
		FormattedError: fmt.Errorf(format, inputs...),
	}
	b.DoReporting()

	return fmt.Errorf(format, inputs...)
}

func (b BugFixesLog) Info(inputs ...interface{}) {
	b.Infof("%v", inputs...)
}
func Info(inputs ...interface{}) {
	Infof("%v", inputs...)
}
func (b BugFixesLog) Infof(format string, inputs ...interface{}) {
	b.Level = "info"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.DoReporting()
}
func Infof(format string, inputs ...interface{}) {
	b := BugFixesLog{
		Level:        "info",
		FormattedLog: fmt.Sprintf(format, inputs...),
	}
	b.DoReporting()
}

func (b BugFixesLog) Debug(inputs ...interface{}) {
	b.Debugf("%v", inputs...)
}
func Debug(inputs ...interface{}) {
	Debugf("%v", inputs...)
}
func (b BugFixesLog) Debugf(format string, inputs ...interface{}) {
	b.Level = "debug"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.Stack = debug.Stack()
	b.DoReporting()
}
func Debugf(format string, inputs ...interface{}) {
	b := BugFixesLog{
		Level:        "debug",
		FormattedLog: fmt.Sprintf(format, inputs...),
		Stack:        debug.Stack(),
	}
	b.DoReporting()
}

func (b BugFixesLog) Log(inputs ...interface{}) {
	b.Logf("%v", inputs...)
}
func Log(inputs ...interface{}) {
	Logf("%v", inputs...)
}
func (b BugFixesLog) Logf(format string, inputs ...interface{}) {
	b.Level = "log"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.DoReporting()
}
func Logf(format string, inputs ...interface{}) {
	b := BugFixesLog{
		Level:        "log",
		FormattedLog: fmt.Sprintf(format, inputs...),
	}
	b.DoReporting()
}

func (b BugFixesLog) Warn(inputs ...interface{}) {
	b.Warnf("%v", inputs...)
}
func Warn(inputs ...interface{}) {
	Warnf("%v", inputs...)
}
func (b BugFixesLog) Warnf(format string, inputs ...interface{}) {
	b.Level = "warn"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.DoReporting()
}
func Warnf(format string, inputs ...interface{}) {
	b := BugFixesLog{
		Level:        "warn",
		FormattedLog: fmt.Sprintf(format, inputs...),
	}
	b.DoReporting()
}
