package logs

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
)

func Errorf(format string, inputs ...interface{}) error {
	go func() {
		_, file, line, _ := runtime.Caller(-1)
		b := BugFixesLog{
			Level:      "error",
			File:       file,
			LineNumber: line,
			Line:       strconv.Itoa(line),
			Log:        fmt.Sprintf(format, inputs...),
		}
		b.DoReporting()
	}()

	return fmt.Errorf(format, inputs...)
}

func Infof(format string, inputs ...interface{}) {
	go func() {
		_, file, line, _ := runtime.Caller(-1)
		b := BugFixesLog{
			Level:      "info",
			File:       file,
			LineNumber: line,
			Line:       strconv.Itoa(line),
			Log:        fmt.Sprintf(format, inputs...),
		}
		b.DoReporting()
	}()
}

func Debugf(format string, inputs ...interface{}) {
	go func() {
		_, file, line, _ := runtime.Caller(-1)
		b := BugFixesLog{
			Level:      "debug",
			File:       file,
			LineNumber: line,
			Line:       strconv.Itoa(line),
			Log:        fmt.Sprintf(format, inputs...),
			Stack:      debug.Stack(),
		}
		b.DoReporting()
	}()
}

func Logf(format string, inputs ...interface{}) {
	go func() {
		_, file, line, _ := runtime.Caller(-1)
		b := BugFixesLog{
			Level:      "log",
			File:       file,
			LineNumber: line,
			Line:       strconv.Itoa(line),
			Log:        fmt.Sprintf(format, inputs...),
		}
		b.DoReporting()
	}()
}
