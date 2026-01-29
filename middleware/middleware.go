package middleware

import (
	"net/http"
)

type System struct {
	// Bugfixes
	AgentID string
	Secret  string

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

func (s *System) AddMiddleware(middlwares ...func(handler http.Handler) http.Handler) {
	s.Middlewares = append(s.Middlewares, middlwares...)
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
	s := NewDefaultMiddleware(next)

	return s.Handler(next)
}

func NewDefaultMiddleware(next http.Handler) *System {
	s := NewMiddleware()
	s.AddMiddleware(Logger)
	s.AddMiddleware(RequestID)
	s.AddMiddleware(Recoverer)
	s.AddMiddleware(LowerCaseHeaders)

	return s
}
