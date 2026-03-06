package middleware

import (
	"net/http"

	bugfixes "github.com/bugfixes/go-bugfixes"
)

type System struct {
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
}

func NewMiddleware() *System {
	return &System{}
}

func (s *System) SetupBugfixes(id, secret string) {
	s.AgentID = id
	s.Secret = secret
}

func (s *System) SetConfig(cfg bugfixes.Config) {
	s.Config = &cfg
}

func (s *System) AddMiddleware(middlewares ...func(handler http.Handler) http.Handler) {
	s.Middlewares = append(s.Middlewares, middlewares...)
}

func (s *System) Handler(h http.Handler) http.Handler {
	if len(s.Middlewares) == 0 {
		return h
	}

	for _, middleware := range s.Middlewares {
		h = middleware(h)
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
