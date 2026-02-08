package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bugfixes/go-bugfixes/middleware"
	"github.com/stretchr/testify/assert"
)

func TestLCHeaders_AddsLowercaseKeys(t *testing.T) {
	var capturedHeaders http.Header
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header
	})

	s := middleware.NewMiddleware()
	handler := s.LCHeaders(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Custom-Header", "value1")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Lowercase key should exist
	assert.Equal(t, "value1", capturedHeaders.Get("x-custom-header"))
	// Original canonical key should still exist
	assert.Equal(t, "value1", capturedHeaders.Get("X-Custom-Header"))
}

func TestLCHeaders_PreservesOriginalHeaders(t *testing.T) {
	var capturedHeaders http.Header
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header
	})

	s := middleware.NewMiddleware()
	handler := s.LCHeaders(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "text/html")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "application/json", capturedHeaders.Get("Accept"))
	assert.Equal(t, "application/json", capturedHeaders.Get("accept"))
	assert.Equal(t, "text/html", capturedHeaders.Get("Content-Type"))
	assert.Equal(t, "text/html", capturedHeaders.Get("content-type"))
}

func TestLCHeaders_MultipleValuesPreserved(t *testing.T) {
	var capturedHeaders http.Header
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header
	})

	s := middleware.NewMiddleware()
	handler := s.LCHeaders(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Multi", "val1")
	req.Header.Add("X-Multi", "val2")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	vals := capturedHeaders.Values("x-multi")
	assert.Len(t, vals, 2)
	assert.Contains(t, vals, "val1")
	assert.Contains(t, vals, "val2")
}

func TestLCHeaders_NoHeaders_NoError(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s := middleware.NewMiddleware()
	handler := s.LCHeaders(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// Clear all default headers
	req.Header = http.Header{}
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestLowerCaseHeaders_FunctionWrapper(t *testing.T) {
	var capturedHeaders http.Header
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header
	})

	handler := middleware.LowerCaseHeaders(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Test", "hello")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "hello", capturedHeaders.Get("x-test"))
	assert.Equal(t, "hello", capturedHeaders.Get("X-Test"))
}
