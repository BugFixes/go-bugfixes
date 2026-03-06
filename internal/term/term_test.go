package term

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorCodes_ValidANSI(t *testing.T) {
	tests := []struct {
		name  string
		color []byte
		want  string
	}{
		{"NRed", NRed, "\033[31m"},
		{"NGreen", NGreen, "\033[32m"},
		{"NYellow", NYellow, "\033[33m"},
		{"NCyan", NCyan, "\033[36m"},
		{"BRed", BRed, "\033[31;1m"},
		{"BGreen", BGreen, "\033[32;1m"},
		{"BYellow", BYellow, "\033[33;1m"},
		{"BBlue", BBlue, "\033[34;1m"},
		{"BMagenta", BMagenta, "\033[35;1m"},
		{"BCyan", BCyan, "\033[36;1m"},
		{"BWhite", BWhite, "\033[37;1m"},
		{"Reset", Reset, "\033[0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, string(tt.color))
		})
	}
}

func TestCW_NoColor(t *testing.T) {
	var buf bytes.Buffer
	CW(&buf, false, BRed, "hello %s", "world")
	assert.Equal(t, "hello world", buf.String())
}

func TestCW_WithColor_NonTTY(t *testing.T) {
	// In test environment, IsTTY is typically false
	origTTY := IsTTY
	IsTTY = false
	defer func() { IsTTY = origTTY }()

	var buf bytes.Buffer
	CW(&buf, true, BRed, "colored")
	assert.Equal(t, "colored", buf.String())
}

func TestCW_WithColor_TTY(t *testing.T) {
	origTTY := IsTTY
	IsTTY = true
	defer func() { IsTTY = origTTY }()

	var buf bytes.Buffer
	CW(&buf, true, BRed, "colored")
	assert.Equal(t, "\033[31;1mcolored\033[0m", buf.String())
}

func TestCW_UseColorFalse_TTY(t *testing.T) {
	origTTY := IsTTY
	IsTTY = true
	defer func() { IsTTY = origTTY }()

	var buf bytes.Buffer
	CW(&buf, false, BRed, "plain")
	assert.Equal(t, "plain", buf.String())
}

func TestCW_FormatArgs(t *testing.T) {
	var buf bytes.Buffer
	CW(&buf, false, BRed, "count: %d, name: %s", 42, "test")
	assert.Equal(t, "count: 42, name: test", buf.String())
}
