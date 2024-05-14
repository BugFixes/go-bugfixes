package logs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/stretchr/testify/assert"
)

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

func Test_Error(t *testing.T) {
	tests := []struct {
		name   string
		inputs error
		expect error
	}{
		{
			name:   "simple",
			inputs: errors.New("simple"),
			expect: fmt.Errorf("%v", errors.New("simple")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := logs.Error(test.inputs)

			if passed := assert.Equal(t, test.expect, response); !passed {
				t.Errorf("%s type failed: %T, expect %T", test.name, response.Error(), test.expect.Error())
			}
		})
	}
}

func Test_Info(t *testing.T) {
	tests := []struct {
		name   string
		inputs error
		expect string
	}{
		{
			name:   "simple",
			inputs: errors.New("simple"),
			expect: fmt.Sprint("Info: simple"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := logs.Info(test.inputs)

			if passed := assert.Equal(t, test.expect, response); !passed {
				t.Errorf("%s type failed: %T, expect %T", test.name, response, test.expect)
			}
		})
	}
}

func Test_Debug(t *testing.T) {
	tests := []struct {
		name   string
		inputs error
		expect string
	}{
		{
			name:   "simple",
			inputs: errors.New("simple"),
			expect: fmt.Sprint("Debug: simple"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := logs.Debug(test.inputs)

			if passed := assert.Equal(t, test.expect, response); !passed {
				t.Errorf("%s type failed: %T, expect %T", test.name, response, test.expect)
			}
		})
	}
}

func Test_Log(t *testing.T) {
	tests := []struct {
		name   string
		inputs error
		expect string
	}{
		{
			name:   "simple",
			inputs: errors.New("simple"),
			expect: fmt.Sprint("Log: simple"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := logs.Log(test.inputs)

			if passed := assert.Equal(t, test.expect, response); !passed {
				t.Errorf("%s type failed: %T, expect %T", test.name, response, test.expect)
			}
		})
	}
}

func Test_Warn(t *testing.T) {
	tests := []struct {
		name   string
		inputs error
		expect string
	}{
		{
			name:   "simple",
			inputs: errors.New("simple"),
			expect: fmt.Sprint("Warn: simple"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := logs.Warn(test.inputs)

			if passed := assert.Equal(t, test.expect, response); !passed {
				t.Errorf("%s type failed: %T, expect %T", test.name, response, test.expect)
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
