package middleware

import (
	"net/http"
)

func BugFixes(next http.Handler) http.Handler {
	handler := RequestID(next)
	handler = Logger(handler)
	return Recoverer(handler)
}
