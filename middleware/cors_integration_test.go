package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bugfixes/go-bugfixes/middleware"
	"github.com/stretchr/testify/assert"
)

// TestCORS_UserSetup mirrors the exact setup from the user's project
func TestCORS_UserSetup(t *testing.T) {
	mw := middleware.NewMiddleware()
	mw.AddMiddleware(middleware.SetupLogger(middleware.Error).Logger)
	mw.AddMiddleware(middleware.RequestID)
	mw.AddMiddleware(middleware.Recoverer)
	mw.AddMiddleware(mw.CORS)
	mw.AddMiddleware(middleware.LowerCaseHeaders)
	mw.AddAllowedMethods(http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodDelete, http.MethodPut)
	mw.AddAllowedOrigins("https://docs.policy.keloran.dev", "http://localhost:4321")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	wrapped := mw.Handler(handler)

	t.Run("OPTIONS preflight returns 200 with CORS headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/endpoint", nil)
		req.Header.Set("Origin", "https://docs.policy.keloran.dev")
		req.Header.Set("Access-Control-Request-Method", "POST")
		rr := httptest.NewRecorder()

		wrapped.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "https://docs.policy.keloran.dev", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, rr.Header().Get("Access-Control-Allow-Methods"), "POST")
		assert.Equal(t, "Origin", rr.Header().Get("Vary"))
	})

	t.Run("POST with allowed origin gets CORS headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/endpoint", nil)
		req.Header.Set("Origin", "https://docs.policy.keloran.dev")
		rr := httptest.NewRecorder()

		wrapped.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "https://docs.policy.keloran.dev", rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("POST with disallowed origin gets 403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/endpoint", nil)
		req.Header.Set("Origin", "https://evil.com")
		rr := httptest.NewRecorder()

		wrapped.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("localhost origin allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/endpoint", nil)
		req.Header.Set("Origin", "http://localhost:4321")
		rr := httptest.NewRecorder()

		wrapped.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "http://localhost:4321", rr.Header().Get("Access-Control-Allow-Origin"))
	})
}

// TestCORS_405_WhenRouterHandlesMethodRouting demonstrates that if a router
// does method-based routing, OPTIONS requests get 405 from the ROUTER, not
// from the CORS middleware. The fix is to wrap the router with middleware.
func TestCORS_405_WhenRouterHandlesMethodRouting(t *testing.T) {
	mw := middleware.NewMiddleware()
	mw.AddMiddleware(mw.CORS)
	mw.AddAllowedOrigins("https://example.com")
	mw.AddAllowedMethods("POST")

	// Simulate a method-routing mux (like Go 1.22+ ServeMux)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("middleware wraps mux - OPTIONS works", func(t *testing.T) {
		// CORRECT: middleware wraps the entire mux
		wrapped := mw.Handler(mux)

		req := httptest.NewRequest(http.MethodOptions, "/api/endpoint", nil)
		req.Header.Set("Origin", "https://example.com")
		rr := httptest.NewRecorder()

		wrapped.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "CORS middleware should intercept OPTIONS before the mux")
		assert.Equal(t, "https://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("mux without middleware - OPTIONS fails", func(t *testing.T) {
		// WITHOUT middleware: mux rejects OPTIONS (no route registered for it)
		req := httptest.NewRequest(http.MethodOptions, "/api/endpoint", nil)
		req.Header.Set("Origin", "https://example.com")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		assert.NotEqual(t, http.StatusOK, rr.Code, "mux rejects OPTIONS when no route is registered for it")
		assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"), "no CORS headers without middleware")
	})
}
