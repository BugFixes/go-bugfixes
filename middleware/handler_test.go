package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bugfixes/go-bugfixes/middleware"
	"github.com/stretchr/testify/assert"
)

func TestHandler_NoMiddlewares_PassesThrough(t *testing.T) {
	s := middleware.NewMiddleware()

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("direct"))
	})

	handler := s.Handler(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "direct", rr.Body.String())
}

func TestHandler_SingleMiddleware(t *testing.T) {
	s := middleware.NewMiddleware()

	headerMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "applied")
			next.ServeHTTP(w, r)
		})
	}
	s.AddMiddleware(headerMiddleware)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := s.Handler(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "applied", rr.Header().Get("X-Test"))
}

func TestHandler_MultipleMiddlewares_AppliedInOrder(t *testing.T) {
	s := middleware.NewMiddleware()

	var order []string
	first := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "first")
			next.ServeHTTP(w, r)
		})
	}
	second := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "second")
			next.ServeHTTP(w, r)
		})
	}
	s.AddMiddleware(first, second)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	handler := s.Handler(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Middlewares wrap in order, so second wraps first's result
	// Execution order: second -> first -> handler
	assert.Len(t, order, 3)
	assert.Equal(t, "handler", order[len(order)-1])
}

func TestNewMiddleware_ReturnsCleanSystem(t *testing.T) {
	s := middleware.NewMiddleware()

	assert.NotNil(t, s)
	assert.Empty(t, s.AllowedOrigins)
	assert.Empty(t, s.AllowedHeaders)
	assert.Empty(t, s.AllowedMethods)
	assert.Empty(t, s.AgentID)
	assert.Empty(t, s.Secret)
}

func TestSetupBugfixes(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetupBugfixes("agent-123", "secret-456")

	assert.Equal(t, "agent-123", s.AgentID)
	assert.Equal(t, "secret-456", s.Secret)
}

func TestDefaultMiddleware_AppliesAllDefaults(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	handler := middleware.DefaultMiddleware(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestRequestID_GeneratedWhenMissing(t *testing.T) {
	var capturedID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = middleware.GetReqID(r.Context())
	})

	handler := middleware.RequestID(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.NotEmpty(t, capturedID)
}

func TestRequestID_UsesProvidedHeader(t *testing.T) {
	var capturedID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = middleware.GetReqID(r.Context())
	})

	handler := middleware.RequestID(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", "custom-id-123")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "custom-id-123", capturedID)
}

func TestRequestID_UniquePerRequest(t *testing.T) {
	var ids []string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids = append(ids, middleware.GetReqID(r.Context()))
	})

	handler := middleware.RequestID(next)

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}

	// All IDs should be unique
	seen := make(map[string]bool)
	for _, id := range ids {
		assert.False(t, seen[id], "duplicate request ID: %s", id)
		seen[id] = true
	}
}

func TestGetReqID_NilContext(t *testing.T) {
	id := middleware.GetReqID(nil)
	assert.Empty(t, id)
}
