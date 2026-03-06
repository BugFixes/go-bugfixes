package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrettyStack_Parse_ContainsPanicValue(t *testing.T) {
	// Simulate a stack trace from runtime
	fakeStack := []byte(`goroutine 1 [running]:
runtime/debug.Stack()
	/usr/local/go/src/runtime/debug/stack.go:24 +0x5e
panic(0x1234, 0x5678)
main.doSomething()
	/app/main.go:42 +0x1a
main.main()
	/app/main.go:10 +0x25
`)

	s := prettyStack{}
	out, err := s.parse(fakeStack, "test panic value")
	require.NoError(t, err)
	assert.Contains(t, string(out), "test panic value")
}

func TestPrettyStack_Parse_EmptyStack(t *testing.T) {
	s := prettyStack{}
	out, err := s.parse([]byte{}, "empty")
	// Should not error, just produce minimal output
	require.NoError(t, err)
	assert.Contains(t, string(out), "empty")
}

func TestPrettyStack_DecorateLine_SourceLine(t *testing.T) {
	s := prettyStack{}
	line := "/app/main.go:42 +0x1a"
	result, err := s.decorateLine(line, false, 0)
	require.NoError(t, err)
	assert.Contains(t, result, "main.go")
	assert.Contains(t, result, ":42")
}

func TestPrettyStack_DecorateLine_FuncCallLine(t *testing.T) {
	s := prettyStack{}
	line := "main.doSomething()"
	result, err := s.decorateLine(line, false, 0)
	require.NoError(t, err)
	assert.Contains(t, result, "doSomething")
}

func TestPrettyStack_DecorateLine_PlainLine(t *testing.T) {
	s := prettyStack{}
	line := "goroutine 1 [running]:"
	// After TrimSpace, this doesn't match source or func patterns, gets default formatting
	result, err := s.decorateLine(line, false, 2)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestPrettyStack_DecorateSourceLine_NotSourceLine(t *testing.T) {
	s := prettyStack{}
	_, err := s.decorateSourceLine("not a source line", false, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a source line")
}

func TestPrettyStack_DecorateFuncCallLine_NotFuncLine(t *testing.T) {
	s := prettyStack{}
	_, err := s.decorateFuncCallLine("no parens here", false, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a func call line")
}

func TestPrettyStack_DecorateSourceLine_HighlightsFirstLine(t *testing.T) {
	s := prettyStack{}
	line := "/app/pkg/handler.go:99 +0x1a"
	result, err := s.decorateSourceLine(line, false, 1)
	require.NoError(t, err)
	assert.Contains(t, result, "handler.go")
	assert.Contains(t, result, ":99")
}

func TestPrettyStack_DecorateFuncCallLine_WithPackage(t *testing.T) {
	s := prettyStack{}
	line := "github.com/example/pkg.Handler()"
	result, err := s.decorateFuncCallLine(line, false, 0)
	require.NoError(t, err)
	assert.Contains(t, result, "Handler")
}

func TestPrettyStack_DecorateFuncCallLine_SimpleFunc(t *testing.T) {
	s := prettyStack{}
	line := "main.run()"
	result, err := s.decorateFuncCallLine(line, false, 1)
	require.NoError(t, err)
	assert.Contains(t, result, "run")
}

func TestPrintPrettyStack_DoesNotPanic(t *testing.T) {
	// PrintPrettyStack should never panic regardless of input
	assert.NotPanics(t, func() {
		PrintPrettyStack("test error")
	})
	assert.NotPanics(t, func() {
		PrintPrettyStack(42)
	})
	assert.NotPanics(t, func() {
		PrintPrettyStack(nil)
	})
}
