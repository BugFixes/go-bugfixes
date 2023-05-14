package logs

import (
	"fmt"
	"runtime/debug"
	"strings"
)

//func SkipLog(level string) bool {
//	type LogLevel struct {
//		String string
//		Int    int
//	}
//	logLevels := []LogLevel{
//		{"log", 0},
//		{"debug", 1},
//		{"info", 2},
//		{"warn", 3},
//		{"error", 4},
//		{"fatal", 5},
//		{"unknown", 9},
//	}
//
//	// default log level
//	logLevel := 0
//	appLogLevel := 0
//
//	logLevelAppRead := os.Getenv("LOG_LEVEL")
//	if logLevelAppRead != "" {
//		i, err := strconv.Atoi(logLevelAppRead)
//		if err != nil {
//			appLogLevel = 0
//		} else {
//			appLogLevel = i
//		}
//	}
//	i, err := strconv.Atoi(level)
//	if err != nil {
//		logLevel = 0
//	} else {
//		logLevel = i
//	}
//
//	for _, logLevel := range logLevels {
//		if logLevel.String == level {
//			return false
//		}
//		if logLevel.Int == appLogLevel {
//			return false
//		}
//	}
//
//	if logLevel >= appLogLevel {
//		return false
//	}
//
//	return true
//}

func Local(skipDepthOverride ...int) BugFixes {
	if len(skipDepthOverride) > 0 && skipDepthOverride[0] != 0 {
		return BugFixes{
			LocalOnly:         true,
			SkipDepthOverride: skipDepthOverride[0],
		}
	}

	return BugFixes{
		LocalOnly: true,
	}
}

// <editor-fold desc="Error">
func (b BugFixes) Error() string {
	return fmt.Sprintf("%s: %s", b.Bug, b.Err.Error())
}

func Error(inputs ...interface{}) error {
	format := strings.Repeat("%w, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return Errorf(format, inputs...)
}
func (b BugFixes) Errorf(format string, inputs ...interface{}) error {
	b.Level = "error"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.FormattedError = fmt.Errorf(format, inputs...)
	b.DoReporting()

	return fmt.Errorf(format, inputs...)
}
func Errorf(format string, inputs ...interface{}) error {
	return BugFixes{
		LocalOnly: false,
	}.Errorf(format, inputs...)
}

// </editor-fold>

// <editor-fold desc="Info">
func (b BugFixes) Info(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return b.Infof(format, inputs...)
}
func Info(inputs ...interface{}) {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	Infof(format, inputs...)
}
func (b BugFixes) Infof(format string, inputs ...interface{}) string {
	b.Level = "info"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.DoReporting()

	return fmt.Sprintf("Info: %s", fmt.Sprintf(format, inputs...))
}
func Infof(format string, inputs ...interface{}) string {
	return BugFixes{
		LocalOnly: false,
	}.Infof(format, inputs...)
}

// </editor-fold>

// <editor-fold desc="Debug">
func (b BugFixes) Debug(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return b.Debugf(format, inputs...)
}
func Debug(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return Debugf(format, inputs...)
}
func (b BugFixes) Debugf(format string, inputs ...interface{}) string {
	b.Level = "debug"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.Stack = debug.Stack()
	b.DoReporting()

	return fmt.Sprintf("Debug: %s", fmt.Sprintf(format, inputs...))
}
func Debugf(format string, inputs ...interface{}) string {
	return BugFixes{
		LocalOnly: false,
	}.Debugf(format, inputs...)
}

// </editor-fold>

// <editor-fold desc="Log">
func (b BugFixes) Log(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return b.Logf(format, inputs...)
}
func Log(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return Logf(format, inputs...)
}
func (b BugFixes) Logf(format string, inputs ...interface{}) string {
	b.Level = "log"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.DoReporting()

	return fmt.Sprintf("Log: %s", fmt.Sprintf(format, inputs...))
}
func Logf(format string, inputs ...interface{}) string {
	return BugFixes{
		LocalOnly: false,
	}.Logf(format, inputs...)
}

// </editor-fold>

// <editor-fold desc="Warn">
func (b BugFixes) Warn(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return b.Warnf(format, inputs...)
}
func Warn(inputs ...interface{}) string {
	format := strings.Repeat("%v, ", len(inputs))
	format = strings.TrimRight(format, ", ") // remove trailing comma and space
	return Warnf(format, inputs...)
}
func (b BugFixes) Warnf(format string, inputs ...interface{}) string {
	b.Level = "warn"
	b.FormattedLog = fmt.Sprintf(format, inputs...)
	b.DoReporting()

	return fmt.Sprintf("Warn: %s", fmt.Sprintf(format, inputs...))
}
func Warnf(format string, inputs ...interface{}) string {
	return BugFixes{
		LocalOnly: false,
	}.Warnf(format, inputs...)
}

// </editor-fold>

func (b BugFixes) Fatal(inputs ...interface{}) {
	b.Level = "fatal"
	b.FormattedLog = fmt.Sprintf("%v", inputs...)
	b.DoReporting()
	panic(b)
}
