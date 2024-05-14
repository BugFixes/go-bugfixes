package middleware

import (
  "context"
  "net/http"
)

type System struct {
	Context context.Context

  // Bugfixes
  AgentID string
  Secret string

	// Middlewares to use
	Middlewares []func(handler http.Handler) http.Handler

	// CORS bits
	AllowedOrigins []string
	AllowedHeaders []string
	AllowedMethods []string
}

func NewMiddleware(ctx context.Context) *System {
	return &System{
		Context: ctx,
	}
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

func DefaultMiddlware(next http.Handler) http.Handler {
	s := NewMiddleware(context.Background())
	s.AddMiddleware(Logger)
	s.AddMiddleware(RequestID)
	s.AddMiddleware(Recoverer)

	return s.Handler(next)
}

