package logs_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bugfixes/go-bugfixes/logs"
)

func TestConvertLevelFromString(t *testing.T) {
	tests := []struct {
		name   string
		level  string
		expect int
	}{
		{
			name:   "log",
			level:  "log",
			expect: logs.LevelLog,
		},
		{
			name:   "info",
			level:  "info",
			expect: logs.LevelInfo,
		},
		{
			name:   "error",
			level:  "error",
			expect: logs.LevelError,
		},
		{
			name:   "warning",
			level:  "warn",
			expect: logs.LevelInfo,
		},
		{
			name:   "unknown",
			level:  "mystery",
			expect: logs.LevelUnknown,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := logs.ConvertLevelFromString(test.level)
			if passed := assert.IsType(t, test.expect, response); !passed {
				t.Errorf("failed type test: %T, expect: %T", response, test.expect)
			}
			if passed := assert.Equal(t, test.expect, response); !passed {
				t.Errorf("failed equal test: %+v, expect: %+v", response, test.expect)
			}
		})
	}
}

func TestBugFixes_Errorf(t *testing.T) {
	tests := []struct {
		name   string
		format string
		inputs string
		expect error
	}{
		{
			name:   "simple",
			format: "this is %s",
			inputs: "simple",
			expect: fmt.Errorf("this is simple"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := logs.Local(3).Errorf(test.format, test.inputs)

			if passed := assert.IsType(t, test.expect, response); !passed {
				t.Errorf("%s type failed: %T, expect %T", test.name, response, test.expect)
			}
			if passed := assert.Equal(t, test.expect, response); !passed {
				t.Errorf("%s equal failed: %+v, expect: %+v", test.name, response, test.expect)
			}
		})
	}
}

func TestBugFixes_Infof(t *testing.T) {
	tests := []struct {
		name   string
		format string
		inputs string
		expect string
	}{
		{
			name:   "simple",
			format: "this is %s",
			inputs: "simple",
			expect: "Info: this is simple",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := logs.Local(3).Infof(test.format, test.inputs)

			if passed := assert.IsType(t, test.expect, response); !passed {
				t.Errorf("%s type failed: %T, expect %T", test.name, response, test.expect)
			}
			if passed := assert.Equal(t, test.expect, response); !passed {
				t.Errorf("%s equal failed: %+v, expect: %+v", test.name, response, test.expect)
			}
		})
	}
}

func TestBugFixes_Debugf(t *testing.T) {
	tests := []struct {
		name   string
		format string
		inputs string
		expect string
	}{
		{
			name:   "simple",
			format: "this is %s",
			inputs: "simple",
			expect: "Debug: this is simple",
		},
	}

	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			response := logs.Local(3).Debugf(test.format, test.inputs)

			if passed := assert.IsType(t, test.expect, response); !passed {
				t.Errorf("%s type failed: %T, expect %T", test.name, response, test.expect)
			}
			if passed := assert.Equal(t, test.expect, response); !passed {
				t.Errorf("%s equal failed: %+v, expect: %+v", test.name, response, test.expect)
			}
		})
	}
}

func TestBugFixes_Logf(t *testing.T) {
	tests := []struct {
		name   string
		format string
		inputs string
		expect string
	}{
		{
			name:   "simple",
			format: "this is %s",
			inputs: "simple",
			expect: "Log: this is simple",
		},
	}

	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			response := logs.Local(3).Logf(test.format, test.inputs)

			if passed := assert.IsType(t, test.expect, response); !passed {
				t.Errorf("%s type failed: %T, expect %T", test.name, response, test.expect)
			}
			if passed := assert.Equal(t, test.expect, response); !passed {
				t.Errorf("%s equal failed: %+v, expect: %+v", test.name, response, test.expect)
			}
		})
	}
}

func TestBugFixes_Warnf(t *testing.T) {
	tests := []struct {
		name   string
		format string
		inputs string
		expect string
	}{
		{
			name:   "simple",
			format: "this is %s",
			inputs: "simple",
			expect: "Warn: this is simple",
		},
	}

	for _, test := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			response := logs.Local(3).Warnf(test.format, test.inputs)

			if passed := assert.IsType(t, test.expect, response); !passed {
				t.Errorf("%s type failed: %T, expect %T", test.name, response, test.expect)
			}
			if passed := assert.Equal(t, test.expect, response); !passed {
				t.Errorf("%s equal failed: %+v, expect: %+v", test.name, response, test.expect)
			}
		})
	}
}
