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
		extras := make(http.Header)
		for k, v := range r.Header {
			lcKey := strings.ToLower(k)
			if lcKey != k {
				extras[lcKey] = v
			}
		}
		for k, v := range extras {
			r.Header[k] = v
		}

		next.ServeHTTP(w, r)
	})
}
