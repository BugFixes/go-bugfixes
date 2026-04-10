package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bugfixes/go-bugfixes/middleware"
	"github.com/stretchr/testify/assert"
)

func TestSecure_DefaultConfig_SetsAllHeaders(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "DENY", rr.Header().Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", rr.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "1; mode=block", rr.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "max-age=31536000; includeSubDomains", rr.Header().Get("Strict-Transport-Security"))
	assert.Equal(t, "default-src 'self'", rr.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "strict-origin-when-cross-origin", rr.Header().Get("Referrer-Policy"))
}

func TestSecure_Disabled_SetsNoHeaders(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(false)

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Empty(t, rr.Header().Get("X-Frame-Options"))
	assert.Empty(t, rr.Header().Get("X-Content-Type-Options"))
	assert.Empty(t, rr.Header().Get("Strict-Transport-Security"))
	assert.Empty(t, rr.Header().Get("Content-Security-Policy"))
}

func TestSecure_CustomConfig_OverridesDefaults(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)
	s.AddSecureConfig("example.com", middleware.SecureConfig{
		XFrameOptions:       middleware.StrPtr("SAMEORIGIN"),
		XContentTypeOptions: middleware.BoolPtr(false),
		CSP:                 middleware.StrPtr("default-src 'self'; script-src 'self' 'unsafe-inline'"),
		ReferrerPolicy:      middleware.StrPtr("no-referrer"),
	})

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "SAMEORIGIN", rr.Header().Get("X-Frame-Options"))
	assert.Empty(t, rr.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "default-src 'self'; script-src 'self' 'unsafe-inline'", rr.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "no-referrer", rr.Header().Get("Referrer-Policy"))
}

func TestSecure_DomainSpecificConfigs(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)

	s.AddSecureConfig("example.com", middleware.SecureConfig{
		XFrameOptions: middleware.StrPtr("SAMEORIGIN"),
		CSP:           middleware.StrPtr("default-src 'self'"),
	})

	s.AddSecureConfig("api.example.com", middleware.SecureConfig{
		XFrameOptions: middleware.StrPtr("DENY"),
		CSP:           middleware.StrPtr("default-src 'self'; connect-src 'self' https://api.example.com"),
	})

	handler := s.Secure(newOKHandler())

	tests := []struct {
		name        string
		host        string
		expectedXFO string
		expectedCSP string
	}{
		{
			name:        "example.com gets SAMEORIGIN",
			host:        "example.com",
			expectedXFO: "SAMEORIGIN",
			expectedCSP: "default-src 'self'",
		},
		{
			name:        "api.example.com gets DENY",
			host:        "api.example.com",
			expectedXFO: "DENY",
			expectedCSP: "default-src 'self'; connect-src 'self' https://api.example.com",
		},
		{
			name:        "unknown domain gets defaults",
			host:        "other.com",
			expectedXFO: "DENY",
			expectedCSP: "default-src 'self'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Host = tt.host
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedXFO, rr.Header().Get("X-Frame-Options"))
			assert.Equal(t, tt.expectedCSP, rr.Header().Get("Content-Security-Policy"))
		})
	}
}

func TestSecure_HSTS_AllOptions(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)
	s.AddSecureConfig("example.com", middleware.SecureConfig{
		HSTSEnabled:           middleware.BoolPtr(true),
		HSTSMaxAge:            middleware.DurationPtr(180 * 24 * time.Hour),
		HSTSIncludeSubdomains: middleware.BoolPtr(true),
		HSTSPreload:           middleware.BoolPtr(true),
	})

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	hsts := rr.Header().Get("Strict-Transport-Security")
	assert.Contains(t, hsts, "max-age=15552000")
	assert.Contains(t, hsts, "includeSubDomains")
	assert.Contains(t, hsts, "preload")
}

func TestSecure_HSTS_Disabled(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)
	s.AddSecureConfig("example.com", middleware.SecureConfig{
		HSTSEnabled: middleware.BoolPtr(false),
	})

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Empty(t, rr.Header().Get("Strict-Transport-Security"))
}

func TestSecure_HSTS_NoSubdomains(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)
	s.AddSecureConfig("example.com", middleware.SecureConfig{
		HSTSEnabled:           middleware.BoolPtr(true),
		HSTSMaxAge:            middleware.DurationPtr(365 * 24 * time.Hour),
		HSTSIncludeSubdomains: middleware.BoolPtr(false),
	})

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	hsts := rr.Header().Get("Strict-Transport-Security")
	assert.Equal(t, "max-age=31536000", hsts)
}

func TestSecure_PermissionsPolicy(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)
	s.AddSecureConfig("example.com", middleware.SecureConfig{
		PermissionsPolicy: middleware.StrPtr("geolocation=(self), microphone=()"),
	})

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "geolocation=(self), microphone=()", rr.Header().Get("Permissions-Policy"))
}

func TestSecure_XXSSProtection_Disabled(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)
	s.AddSecureConfig("example.com", middleware.SecureConfig{
		XXSSProtection: middleware.StrPtr(""),
	})

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Empty(t, rr.Header().Get("X-XSS-Protection"))
}

func TestSecure_XFrameOptions_AllowFrom(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)
	s.AddSecureConfig("example.com", middleware.SecureConfig{
		XFrameOptions: middleware.StrPtr("ALLOW-FROM https://trusted.com"),
	})

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "ALLOW-FROM https://trusted.com", rr.Header().Get("X-Frame-Options"))
}

func TestSecure_NextHandlerIsCalled(t *testing.T) {
	called := false
	s := middleware.NewMiddleware()
	s.SetSecure(true)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	handler := s.Secure(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.True(t, called, "next handler should be called")
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestSecure_MultipleSetSecureCalls(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(false)
	s.SetSecure(true)

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.NotEmpty(t, rr.Header().Get("X-Frame-Options"))
}

func TestSecure_AppendSecureConfigs(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)

	s.AddSecureConfig("one.com", middleware.SecureConfig{
		XFrameOptions: middleware.StrPtr("SAMEORIGIN"),
	})

	s.AddSecureConfig("two.com", middleware.SecureConfig{
		XFrameOptions: middleware.StrPtr("DENY"),
	})

	handler := s.Secure(newOKHandler())

	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.Host = "one.com"
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	assert.Equal(t, "SAMEORIGIN", rr1.Header().Get("X-Frame-Options"))

	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.Host = "two.com"
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	assert.Equal(t, "DENY", rr2.Header().Get("X-Frame-Options"))
}

func TestSetSecureConfig_UpdatesExisting(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)

	s.SetSecureConfig("example.com", middleware.SecureConfig{
		XFrameOptions: middleware.StrPtr("SAMEORIGIN"),
	})

	s.SetSecureConfig("example.com", middleware.SecureConfig{
		XFrameOptions: middleware.StrPtr("DENY"),
	})

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "DENY", rr.Header().Get("X-Frame-Options"))
}

func TestSecureConfig_WithHostPort(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)

	s.AddSecureConfig("example.com:8080", middleware.SecureConfig{
		XFrameOptions: middleware.StrPtr("SAMEORIGIN"),
	})

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com:8080"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "SAMEORIGIN", rr.Header().Get("X-Frame-Options"))
}

func TestSecure_EmptyHostFallsBackToDefault(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = ""
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "DENY", rr.Header().Get("X-Frame-Options"))
	assert.Equal(t, "default-src 'self'", rr.Header().Get("Content-Security-Policy"))
}

func TestSecure_CSPWithReportURI(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)
	s.AddSecureConfig("example.com", middleware.SecureConfig{
		CSP: middleware.StrPtr("default-src 'self'; report-uri /csp-report"),
	})

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "default-src 'self'; report-uri /csp-report", rr.Header().Get("Content-Security-Policy"))
}

func TestSecure_ReferrerPolicy_Values(t *testing.T) {
	testCases := []struct {
		policy   string
		expected string
	}{
		{"no-referrer", "no-referrer"},
		{"no-referrer-when-downgrade", "no-referrer-when-downgrade"},
		{"origin", "origin"},
		{"origin-when-cross-origin", "origin-when-cross-origin"},
		{"same-origin", "same-origin"},
		{"strict-origin", "strict-origin"},
		{"strict-origin-when-cross-origin", "strict-origin-when-cross-origin"},
	}

	for _, tc := range testCases {
		t.Run(tc.policy, func(t *testing.T) {
			s := middleware.NewMiddleware()
			s.SetSecure(true)
			s.AddSecureConfig("example.com", middleware.SecureConfig{
				ReferrerPolicy: middleware.StrPtr(tc.policy),
			})

			handler := s.Secure(newOKHandler())
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Host = "example.com"
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expected, rr.Header().Get("Referrer-Policy"))
		})
	}
}

func TestSecure_XContentTypeOptions_NosniffOnly(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)
	s.AddSecureConfig("example.com", middleware.SecureConfig{
		XContentTypeOptions: middleware.BoolPtr(true),
	})

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "nosniff", rr.Header().Get("X-Content-Type-Options"))
}

func TestDefaultSecureConfig_Values(t *testing.T) {
	cfg := middleware.DefaultSecureConfig

	assert.Equal(t, "DENY", *cfg.XFrameOptions)
	assert.True(t, *cfg.XContentTypeOptions)
	assert.Equal(t, "1; mode=block", *cfg.XXSSProtection)
	assert.True(t, *cfg.HSTSEnabled)
	assert.Equal(t, 365*24*time.Hour, *cfg.HSTSMaxAge)
	assert.True(t, *cfg.HSTSIncludeSubdomains)
	assert.False(t, *cfg.HSTSPreload)
	assert.Equal(t, "default-src 'self'", *cfg.CSP)
	assert.Equal(t, "strict-origin-when-cross-origin", *cfg.ReferrerPolicy)
	assert.Nil(t, cfg.PermissionsPolicy)
}

func TestSecure_PartialConfig_MergesWithDefaults(t *testing.T) {
	s := middleware.NewMiddleware()
	s.SetSecure(true)

	s.AddSecureConfig("example.com", middleware.SecureConfig{
		XFrameOptions: middleware.StrPtr("SAMEORIGIN"),
	})

	handler := s.Secure(newOKHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "example.com"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "SAMEORIGIN", rr.Header().Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", rr.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "max-age=31536000; includeSubDomains", rr.Header().Get("Strict-Transport-Security"))
	assert.Equal(t, "default-src 'self'", rr.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "strict-origin-when-cross-origin", rr.Header().Get("Referrer-Policy"))
}
