package middleware

import (
	"net/http"
	"sync"
	"time"

	bugfixes "github.com/bugfixes/go-bugfixes"
)

type SecureConfig struct {
	XFrameOptions         *string
	XContentTypeOptions   *bool
	XXSSProtection        *string
	HSTSEnabled           *bool
	HSTSMaxAge            *time.Duration
	HSTSIncludeSubdomains *bool
	HSTSPreload           *bool
	CSP                   *string
	ReferrerPolicy        *string
	PermissionsPolicy     *string
}

var DefaultSecureConfig = SecureConfig{
	XFrameOptions:         StrPtr("DENY"),
	XContentTypeOptions:   BoolPtr(true),
	XXSSProtection:        StrPtr("1; mode=block"),
	HSTSEnabled:           BoolPtr(true),
	HSTSMaxAge:            DurationPtr(365 * 24 * time.Hour),
	HSTSIncludeSubdomains: BoolPtr(true),
	HSTSPreload:           BoolPtr(false),
	ReferrerPolicy:        StrPtr("strict-origin-when-cross-origin"),
	CSP:                   StrPtr("default-src 'self'"),
}

func StrPtr(s string) *string                    { return &s }
func BoolPtr(b bool) *bool                       { return &b }
func DurationPtr(d time.Duration) *time.Duration { return &d }

type System struct {
	mu sync.RWMutex

	// Bugfixes
	AgentID string
	Secret  string
	Config  *bugfixes.Config

	// Middlewares to use
	Middlewares []func(handler http.Handler) http.Handler

	// CORS bits
	AllowedOrigins []string
	AllowedHeaders []string
	AllowedMethods []string

	// Secure headers
	SecureEnabled bool
	SecureConfigs map[string]SecureConfig
}

func NewMiddleware() *System {
	return &System{}
}

func (s *System) SetupBugfixes(id, secret string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.AgentID = id
	s.Secret = secret
}

func (s *System) SetConfig(cfg bugfixes.Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Config = &cfg
}

func (s *System) AddMiddleware(middlewares ...func(handler http.Handler) http.Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Middlewares = append(s.Middlewares, middlewares...)
}

// Handler applies the registered middlewares to h.
// Middlewares execute in registration order: the first middleware added
// is the outermost handler (runs first).
func (s *System) Handler(h http.Handler) http.Handler {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.Middlewares) == 0 {
		return h
	}

	// Apply in reverse so the first-registered middleware wraps outermost.
	for i := len(s.Middlewares) - 1; i >= 0; i-- {
		h = s.Middlewares[i](h)
	}

	return h
}

func DefaultMiddleware(next http.Handler) http.Handler {
	return NewDefaultMiddleware().Handler(next)
}

func NewDefaultMiddleware() *System {
	s := NewMiddleware()
	s.AddMiddleware(Logger)
	s.AddMiddleware(RequestID)
	s.AddMiddleware(Recoverer)
	s.AddMiddleware(LowerCaseHeaders)

	return s
}
