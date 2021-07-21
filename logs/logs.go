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

// <editor-fold desc="Error">
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
	return BugFixes{}.Errorf(format, inputs...)
}

// </editor-fold>

// <editor-fold desc="Info">
func (b BugFixes) Info(inputs ...interface{}) string {
	return b.Infof("%v", inputs...)
}
func Info(inputs ...interface{}) {
	Infof("%v", inputs...)
}
func (b BugFixes) Infof(format string, inputs ...interface{}) string {
	b.Level = "info"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.DoReporting()

	return fmt.Sprintf("Info: %s", fmt.Sprintf(format, inputs...))
}
func Infof(format string, inputs ...interface{}) string {
	return BugFixes{}.Infof(format, inputs...)
}

// </editor-fold>

// <editor-fold desc="Debug">
func (b BugFixes) Debug(inputs ...interface{}) string {
	return b.Debugf("%v", inputs...)
}
func Debug(inputs ...interface{}) string {
	return Debugf("%v", inputs...)
}
func (b BugFixes) Debugf(format string, inputs ...interface{}) string {
	b.Level = "debug"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.Stack = debug.Stack()
	b.DoReporting()

	return fmt.Sprintf("Debug: %s", fmt.Sprintf(format, inputs...))
}
func Debugf(format string, inputs ...interface{}) string {
	return BugFixes{}.Debugf(format, inputs...)
}

// </editor-fold>

// <editor-fold desc="Log">
func (b BugFixes) Log(inputs ...interface{}) string {
	return b.Logf("%v", inputs...)
}
func Log(inputs ...interface{}) string {
	return Logf("%v", inputs...)
}
func (b BugFixes) Logf(format string, inputs ...interface{}) string {
	b.Level = "log"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.DoReporting()

	return fmt.Sprintf("Log: %s", fmt.Sprintf(format, inputs...))
}
func Logf(format string, inputs ...interface{}) string {
	return BugFixes{}.Logf(format, inputs...)
}

// </editor-fold>

// <editor-fold desc="Warn">
func (b BugFixes) Warn(inputs ...interface{}) string {
	return b.Warnf("%v", inputs...)
}
func Warn(inputs ...interface{}) string {
	return Warnf("%v", inputs...)
}
func (b BugFixes) Warnf(format string, inputs ...interface{}) string {
	b.Level = "warn"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.DoReporting()

	return fmt.Sprintf("Warn: %s", fmt.Sprintf(format, inputs...))
}
func Warnf(format string, inputs ...interface{}) string {
	return BugFixes{}.Warnf(format, inputs...)
}

// </editor-fold>
