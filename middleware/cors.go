package middleware

import (
  "net/http"
  "strings"
)

// Origins
func (s *System) AddAllowedOrigins(origins ...string) {
	s.AllowedOrigins = append(s.AllowedOrigins, origins...)
}

// Headers
func (s *System) AddAllowedHeaders(headers ...string) {
	s.AllowedHeaders = append(s.AllowedHeaders, headers...)
}
func (s *System) getAllowedHeaders() string {
	standardAllowed := []string {
		"Accept",
		"Content-Type",
	}
	
	allowedHeaders := append(standardAllowed, s.AllowedHeaders...)
	return strings.Join(allowedHeaders, ", ")
}

// Methods
func (s *System) AddAllowedMethods(methods ...string) {
	s.AllowedMethods = append(s.AllowedMethods, methods...)
}
func (s *System) getAllowedMethods() string {
	return strings.Join(s.AllowedMethods, ", ")
}

func (s *System) wildcardEnabled() bool {
  for _, origin := range s.AllowedOrigins {
    if origin == "*" {
      return true
    }
  }

  return false
}


func (s *System) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalOrigin := r.Header.Get("Origin")

		// No Origin header means this is not a CORS request, let it through
		if originalOrigin == "" {
			next.ServeHTTP(w, r)
			return
		}

		isAllowed := s.wildcardEnabled()
		for _, origin := range s.AllowedOrigins {
			if origin == originalOrigin {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if s.wildcardEnabled() {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", originalOrigin)
		}
		w.Header().Set("Access-Control-Allow-Methods", s.getAllowedMethods())
		w.Header().Set("Access-Control-Allow-Headers", s.getAllowedHeaders())
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Vary", "Origin")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
