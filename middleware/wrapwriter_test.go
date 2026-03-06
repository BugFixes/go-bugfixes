package middleware_test

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bugfixes/go-bugfixes/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicWriter_DefaultStatusOK(t *testing.T) {
	rr := httptest.NewRecorder()
	w := middleware.NewWrapResponseWriter(rr, 1)

	_, err := w.Write([]byte("hello"))
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Status())
	assert.Equal(t, int64(5), w.BytesWritten())
}

func TestBasicWriter_ExplicitStatus(t *testing.T) {
	rr := httptest.NewRecorder()
	w := middleware.NewWrapResponseWriter(rr, 1)

	w.WriteHeader(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, w.Status())
}

func TestBasicWriter_WriteHeaderIdempotent(t *testing.T) {
	rr := httptest.NewRecorder()
	w := middleware.NewWrapResponseWriter(rr, 1)

	w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusInternalServerError) // should be ignored

	assert.Equal(t, http.StatusCreated, w.Status())
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestBasicWriter_BytesWritten(t *testing.T) {
	rr := httptest.NewRecorder()
	w := middleware.NewWrapResponseWriter(rr, 1)

	_, _ = w.Write([]byte("abc"))
	_, _ = w.Write([]byte("defgh"))

	assert.Equal(t, int64(8), w.BytesWritten())
}

func TestBasicWriter_Tee(t *testing.T) {
	rr := httptest.NewRecorder()
	w := middleware.NewWrapResponseWriter(rr, 1)

	var buf bytes.Buffer
	w.Tee(&buf)

	_, err := w.Write([]byte("tee test"))
	require.NoError(t, err)

	assert.Equal(t, "tee test", rr.Body.String())
	assert.Equal(t, "tee test", buf.String())
}

func TestBasicWriter_Unwrap(t *testing.T) {
	rr := httptest.NewRecorder()
	w := middleware.NewWrapResponseWriter(rr, 1)

	assert.Equal(t, rr, w.Unwrap())
}

func TestBasicWriter_HeaderPassthrough(t *testing.T) {
	rr := httptest.NewRecorder()
	w := middleware.NewWrapResponseWriter(rr, 1)

	w.Header().Set("X-Custom", "value")
	assert.Equal(t, "value", rr.Header().Get("X-Custom"))
}

// mockFlusher implements http.ResponseWriter and http.Flusher.
type mockFlusher struct {
	http.ResponseWriter
	flushed bool
}

func (m *mockFlusher) Flush() { m.flushed = true }

func TestFlushWriter(t *testing.T) {
	rr := httptest.NewRecorder()
	mf := &mockFlusher{ResponseWriter: rr}
	w := middleware.NewWrapResponseWriter(mf, 1)

	_, ok := w.(http.Flusher)
	assert.True(t, ok, "should implement http.Flusher")

	w.(http.Flusher).Flush()
	assert.True(t, mf.flushed)
}

// mockHijacker implements http.ResponseWriter and http.Hijacker.
type mockHijacker struct {
	http.ResponseWriter
	hijacked bool
}

func (m *mockHijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	m.hijacked = true
	return nil, nil, nil
}

func TestHijackWriter(t *testing.T) {
	rr := httptest.NewRecorder()
	mh := &mockHijacker{ResponseWriter: rr}
	w := middleware.NewWrapResponseWriter(mh, 1)

	_, ok := w.(http.Hijacker)
	assert.True(t, ok, "should implement http.Hijacker")

	_, _, err := w.(http.Hijacker).Hijack()
	assert.NoError(t, err)
	assert.True(t, mh.hijacked)
}

// mockFlushHijacker implements Flusher + Hijacker.
type mockFlushHijacker struct {
	http.ResponseWriter
	flushed  bool
	hijacked bool
}

func (m *mockFlushHijacker) Flush()                                        { m.flushed = true }
func (m *mockFlushHijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) { m.hijacked = true; return nil, nil, nil }

func TestFlushHijackWriter(t *testing.T) {
	rr := httptest.NewRecorder()
	mfh := &mockFlushHijacker{ResponseWriter: rr}
	w := middleware.NewWrapResponseWriter(mfh, 1)

	_, isFlusher := w.(http.Flusher)
	_, isHijacker := w.(http.Hijacker)
	assert.True(t, isFlusher, "should implement http.Flusher")
	assert.True(t, isHijacker, "should implement http.Hijacker")

	w.(http.Flusher).Flush()
	assert.True(t, mfh.flushed)

	_, _, err := w.(http.Hijacker).Hijack()
	assert.NoError(t, err)
	assert.True(t, mfh.hijacked)
}

// mockFancyWriter implements Flusher + Hijacker + ReaderFrom.
type mockFancyWriter struct {
	http.ResponseWriter
	flushed  bool
	hijacked bool
}

func (m *mockFancyWriter) Flush()                                        { m.flushed = true }
func (m *mockFancyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) { m.hijacked = true; return nil, nil, nil }
func (m *mockFancyWriter) ReadFrom(r io.Reader) (int64, error)          { return io.Copy(m.ResponseWriter, r) }

func TestHTTPFancyWriter(t *testing.T) {
	rr := httptest.NewRecorder()
	mfw := &mockFancyWriter{ResponseWriter: rr}
	w := middleware.NewWrapResponseWriter(mfw, 1)

	_, isFlusher := w.(http.Flusher)
	_, isHijacker := w.(http.Hijacker)
	_, isReaderFrom := w.(io.ReaderFrom)

	assert.True(t, isFlusher, "should implement http.Flusher")
	assert.True(t, isHijacker, "should implement http.Hijacker")
	assert.True(t, isReaderFrom, "should implement io.ReaderFrom")
}

func TestHTTPFancyWriter_ReadFrom(t *testing.T) {
	rr := httptest.NewRecorder()
	mfw := &mockFancyWriter{ResponseWriter: rr}
	w := middleware.NewWrapResponseWriter(mfw, 1)

	rf := w.(io.ReaderFrom)
	n, err := rf.ReadFrom(bytes.NewReader([]byte("read from test")))
	require.NoError(t, err)
	assert.Equal(t, int64(14), n)
	assert.Equal(t, int64(14), w.BytesWritten())
}

func TestHTTPFancyWriter_ReadFromWithTee(t *testing.T) {
	rr := httptest.NewRecorder()
	mfw := &mockFancyWriter{ResponseWriter: rr}
	w := middleware.NewWrapResponseWriter(mfw, 1)

	var teeBuf bytes.Buffer
	w.Tee(&teeBuf)

	rf := w.(io.ReaderFrom)
	n, err := rf.ReadFrom(bytes.NewReader([]byte("tee read")))
	require.NoError(t, err)
	assert.Equal(t, int64(8), n)
	assert.Equal(t, "tee read", teeBuf.String())
}

// mockHTTP2Writer implements Flusher + Pusher for HTTP/2.
type mockHTTP2Writer struct {
	http.ResponseWriter
	flushed bool
	pushed  bool
}

func (m *mockHTTP2Writer) Flush() { m.flushed = true }
func (m *mockHTTP2Writer) Push(target string, opts *http.PushOptions) error {
	m.pushed = true
	return nil
}

func TestHTTP2FancyWriter(t *testing.T) {
	rr := httptest.NewRecorder()
	mw := &mockHTTP2Writer{ResponseWriter: rr}
	w := middleware.NewWrapResponseWriter(mw, 2)

	_, isFlusher := w.(http.Flusher)
	_, isPusher := w.(http.Pusher)

	assert.True(t, isFlusher, "should implement http.Flusher")
	assert.True(t, isPusher, "should implement http.Pusher")

	w.(http.Flusher).Flush()
	assert.True(t, mw.flushed)

	err := w.(http.Pusher).Push("/resource", nil)
	assert.NoError(t, err)
	assert.True(t, mw.pushed)
}

func TestHTTP2_FallsBackToFlushWriter(t *testing.T) {
	// HTTP/2 without Pusher should get a flushWriter
	rr := httptest.NewRecorder()
	mf := &mockFlusher{ResponseWriter: rr}
	w := middleware.NewWrapResponseWriter(mf, 2)

	_, isFlusher := w.(http.Flusher)
	_, isPusher := w.(http.Pusher)

	assert.True(t, isFlusher)
	assert.False(t, isPusher)
}

func TestProtoMajor1_NoInterfaces_ReturnsBasic(t *testing.T) {
	rr := httptest.NewRecorder()
	w := middleware.NewWrapResponseWriter(rr, 1)

	// httptest.ResponseRecorder implements Flusher, so wrap a plain writer
	plain := &plainWriter{ResponseWriter: rr}
	w2 := middleware.NewWrapResponseWriter(plain, 1)

	_, isFlusher := w2.(http.Flusher)
	assert.False(t, isFlusher, "plain writer should not implement Flusher")

	// But httptest.ResponseRecorder does implement Flusher
	_, isFlusher = w.(http.Flusher)
	assert.True(t, isFlusher)
}

type plainWriter struct {
	http.ResponseWriter
}

func (p *plainWriter) Header() http.Header         { return p.ResponseWriter.Header() }
func (p *plainWriter) Write(b []byte) (int, error)  { return p.ResponseWriter.Write(b) }
func (p *plainWriter) WriteHeader(s int)            { p.ResponseWriter.WriteHeader(s) }
