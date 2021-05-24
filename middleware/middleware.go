package middleware

import (
	"net/http"

	"github.com/bugfixes/go-bugfixes"
)

func Middleware(next http.Handler) http.Handler {
	handler := bugfixes.Logger(next)
	return bugfixes.Recoverer(handler)
}
