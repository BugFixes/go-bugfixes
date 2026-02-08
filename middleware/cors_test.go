package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bugfixes/go-bugfixes/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newOKHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
}

func TestCORS_NoOriginHeader_PassesThrough(t *testing.T) {
	s := middleware.NewMiddleware()
	s.AddAllowedOrigins("https://example.com")

	handler := s.CORS(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
	assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_AllowedOrigin_SetsHeaders(t *testing.T) {
	s := middleware.NewMiddleware()
	s.AddAllowedOrigins("https://example.com")
	s.AddAllowedMethods("GET", "POST")
	s.AddAllowedHeaders("Authorization")

	handler := s.CORS(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://example.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "https://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST", rr.Header().Get("Access-Control-Allow-Methods"))
	assert.Contains(t, rr.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	assert.Contains(t, rr.Header().Get("Access-Control-Allow-Headers"), "Accept")
	assert.Contains(t, rr.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	assert.Equal(t, "86400", rr.Header().Get("Access-Control-Max-Age"))
	assert.Equal(t, "Origin", rr.Header().Get("Vary"))
}

func TestCORS_DisallowedOrigin_Returns403(t *testing.T) {
	s := middleware.NewMiddleware()
	s.AddAllowedOrigins("https://example.com")

	handler := s.CORS(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://evil.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_Wildcard_SetsStar(t *testing.T) {
	s := middleware.NewMiddleware()
	s.AddAllowedOrigins("*")

	handler := s.CORS(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://anything.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_WildcardAmongOthers_SetsStar(t *testing.T) {
	s := middleware.NewMiddleware()
	s.AddAllowedOrigins("https://example.com", "*")

	handler := s.CORS(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://other.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_OptionsPreflight_Returns200(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	s := middleware.NewMiddleware()
	s.AddAllowedOrigins("https://example.com")
	s.AddAllowedMethods("GET", "POST", "PUT")

	handler := s.CORS(next)
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "https://example.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.False(t, called, "next handler should not be called for preflight")
	assert.Equal(t, "https://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_MultipleOrigins_MatchesCorrectOne(t *testing.T) {
	s := middleware.NewMiddleware()
	s.AddAllowedOrigins("https://one.com", "https://two.com", "https://three.com")

	handler := s.CORS(newOKHandler())

	tests := []struct {
		name   string
		origin string
		code   int
	}{
		{"first origin", "https://one.com", http.StatusOK},
		{"second origin", "https://two.com", http.StatusOK},
		{"third origin", "https://three.com", http.StatusOK},
		{"unknown origin", "https://four.com", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Origin", tt.origin)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.code, rr.Code)
			if tt.code == http.StatusOK {
				assert.Equal(t, tt.origin, rr.Header().Get("Access-Control-Allow-Origin"))
			}
		})
	}
}

func TestCORS_NoAllowedOrigins_BlocksAll(t *testing.T) {
	s := middleware.NewMiddleware()

	handler := s.CORS(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://example.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestCORS_NonOptionsMethod_CallsNextHandler(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	s := middleware.NewMiddleware()
	s.AddAllowedOrigins("https://example.com")

	handler := s.CORS(next)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Origin", "https://example.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.True(t, called, "next handler should be called for non-OPTIONS")
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestAddAllowedOrigins_AppendsMultiple(t *testing.T) {
	s := middleware.NewMiddleware()
	s.AddAllowedOrigins("https://one.com")
	s.AddAllowedOrigins("https://two.com", "https://three.com")

	require.Len(t, s.AllowedOrigins, 3)
	assert.Equal(t, "https://one.com", s.AllowedOrigins[0])
	assert.Equal(t, "https://two.com", s.AllowedOrigins[1])
	assert.Equal(t, "https://three.com", s.AllowedOrigins[2])
}

func TestAddAllowedMethods_AppendsMultiple(t *testing.T) {
	s := middleware.NewMiddleware()
	s.AddAllowedMethods("GET")
	s.AddAllowedMethods("POST", "PUT")

	require.Len(t, s.AllowedMethods, 3)
}

func TestAddAllowedHeaders_AppendsMultiple(t *testing.T) {
	s := middleware.NewMiddleware()
	s.AddAllowedHeaders("Authorization")
	s.AddAllowedHeaders("X-Custom", "X-Other")

	require.Len(t, s.AllowedHeaders, 3)
}
