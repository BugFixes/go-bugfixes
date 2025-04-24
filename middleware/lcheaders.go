package middleware

import (
	"net/http"
	"strings"
)

func LowerCaseHeaders(next http.Handler) http.Handler {
	s := &System{}
	return s.LCHeaders(next)
}

func (s *System) LCHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Iterate over the request headers
		for k, v := range r.Header {
			// Create a lowercase version of the key
			lcKey := strings.ToLower(k)

			// Add the lowercase key with the same values
			// This won't overwrite the original header
			r.Header[lcKey] = v
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
