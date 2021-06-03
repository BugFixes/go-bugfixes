package logs

import (
	"fmt"
	"runtime/debug"
)

func Error(inputs ...interface{}) error {
	return Errorf("%w", inputs...)
}
func Errorf(format string, inputs ...interface{}) error {
	b := BugFixesLog{
		Level: "error",
		Log:   fmt.Sprintf(format, inputs...),
		Error: fmt.Errorf(format, inputs...),
	}
	b.DoReporting()

	return fmt.Errorf(format, inputs...)
}

func Info(inputs ...interface{}) {
	Infof("%v", inputs...)
}
func Infof(format string, inputs ...interface{}) {
	b := BugFixesLog{
		Level: "info",
		Log:   fmt.Sprintf(format, inputs...),
	}
	b.DoReporting()
}

func Debug(inputs ...interface{}) {
	Debugf("%v", inputs...)
}
func Debugf(format string, inputs ...interface{}) {
	b := BugFixesLog{
		Level: "debug",
		Log:   fmt.Sprintf(format, inputs...),
		Stack: debug.Stack(),
	}
	b.DoReporting()
}

func Log(inputs ...interface{}) {
	Logf("%v", inputs...)
}
func Logf(format string, inputs ...interface{}) {
	b := BugFixesLog{
		Level: "log",
		Log:   fmt.Sprintf(format, inputs...),
	}
	b.DoReporting()
}

func Warn(inputs ...interface{}) {
	Warnf("%v", inputs...)
}
func Warnf(format string, inputs ...interface{}) {
	b := BugFixesLog{
		Level: "warn",
		Log:   fmt.Sprintf(format, inputs...),
	}
	b.DoReporting()
}
