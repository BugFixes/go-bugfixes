package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bugfixes/go-bugfixes/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRecoverer_PanicReturns500(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	handler := middleware.Recoverer(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	// Should not panic
	assert.NotPanics(t, func() {
		handler.ServeHTTP(rr, req)
	})

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRecoverer_NoPanic_PassesThrough(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	handler := middleware.Recoverer(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestRecoverer_PanicWithString(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic string")
	})

	handler := middleware.Recoverer(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		handler.ServeHTTP(rr, req)
	})
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRecoverer_PanicWithError(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(assert.AnError)
	})

	handler := middleware.Recoverer(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		handler.ServeHTTP(rr, req)
	})
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRecoverer_PanicWithInt(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(42)
	})

	handler := middleware.Recoverer(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		handler.ServeHTTP(rr, req)
	})
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRecoverer_ErrAbortHandler_SilentlyIgnored(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(http.ErrAbortHandler)
	})

	handler := middleware.Recoverer(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	// ErrAbortHandler is recovered but silently ignored — no 500 written
	assert.NotPanics(t, func() {
		handler.ServeHTTP(rr, req)
	})

	assert.Equal(t, http.StatusOK, rr.Code, "no status should be written for ErrAbortHandler")
}

func TestRecoverer_SystemMethod(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("system recoverer test")
	})

	s := middleware.NewMiddleware()
	handler := s.Recoverer(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		handler.ServeHTTP(rr, req)
	})
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
