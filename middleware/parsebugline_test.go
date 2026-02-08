package middleware_test

import (
	"testing"

	"github.com/bugfixes/go-bugfixes/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBugLine_Valid(t *testing.T) {
	file, lne, line, err := middleware.ParseBugLine("main.go:42 +0x1a")

	require.NoError(t, err)
	assert.Equal(t, "main.go", file)
	assert.Equal(t, "42", lne)
	assert.Equal(t, 42, line)
}

func TestParseBugLine_DeepPath(t *testing.T) {
	file, lne, line, err := middleware.ParseBugLine("pkg/handler.go:100 +0xff")

	require.NoError(t, err)
	assert.Equal(t, "pkg/handler.go", file)
	assert.Equal(t, "100", lne)
	assert.Equal(t, 100, line)
}

func TestParseBugLine_MissingColon_ReturnsError(t *testing.T) {
	_, _, _, err := middleware.ParseBugLine("main.go 42")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "':'")
}

func TestParseBugLine_MissingSpace_ReturnsError(t *testing.T) {
	_, _, _, err := middleware.ParseBugLine("main.go:42")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "' '")
}

func TestParseBugLine_EmptyString_ReturnsError(t *testing.T) {
	_, _, _, err := middleware.ParseBugLine("")

	require.Error(t, err)
}

func TestParseBugLine_NonNumericLine_ReturnsError(t *testing.T) {
	file, lne, line, err := middleware.ParseBugLine("main.go:abc +0x1a")

	require.Error(t, err)
	assert.Equal(t, "main.go", file)
	assert.Equal(t, "abc", lne)
	assert.Equal(t, 0, line)
	assert.Contains(t, err.Error(), "convert line number")
}

func TestParseBugLine_ColonOnly_ReturnsError(t *testing.T) {
	_, _, _, err := middleware.ParseBugLine(":")

	require.Error(t, err)
}

func TestParseBugLine_LineNumberZero(t *testing.T) {
	file, lne, line, err := middleware.ParseBugLine("test.go:0 +0x00")

	require.NoError(t, err)
	assert.Equal(t, "test.go", file)
	assert.Equal(t, "0", lne)
	assert.Equal(t, 0, line)
}

func TestParseBugLine_LargeLineNumber(t *testing.T) {
	file, lne, line, err := middleware.ParseBugLine("big.go:99999 +0x00")

	require.NoError(t, err)
	assert.Equal(t, "big.go", file)
	assert.Equal(t, "99999", lne)
	assert.Equal(t, 99999, line)
}
