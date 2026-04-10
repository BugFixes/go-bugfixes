# go-bugfixes

`go-bugfixes` is a small Go library that provides:

- application logging helpers in `logs`
- HTTP middleware in `middleware`

## Requirements

The library targets Go 1.26+.

To send logs or panic reports to Bugfixes, set:

- `BUGFIXES_AGENT_KEY`
- `BUGFIXES_AGENT_SECRET`

Optional environment variables:

- `BUGFIXES_LOCAL_ONLY=true` keeps reporting local
- `BUGFIXES_LOG_LEVEL` sets the minimum remote reporting level
- `BUGFIXES_SERVER` overrides the default API endpoint

## Install

```bash
go get github.com/bugfixes/go-bugfixes
```

## Logging

```go
package main

import "github.com/bugfixes/go-bugfixes/logs"

func main() {
	logger := logs.Local()
	_ = logger.Infof("server started on %s", ":8080")
}
```

Package-level helpers are also available:

```go
package main

import "github.com/bugfixes/go-bugfixes/logs"

func main() {
	_ = logs.Warn("cache is warming")
}
```

## Middleware

The middleware package is router-agnostic and works with standard `net/http` middleware chains.

```go
package main

import (
	"net/http"

	bugfixes "github.com/bugfixes/go-bugfixes/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	handler := bugfixes.DefaultMiddleware(mux)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
```

For a manual chain, add the middleware you need:

```go
system := bugfixes.NewMiddleware()
system.AddMiddleware(
	bugfixes.Logger,
	bugfixes.RequestID,
	bugfixes.Recoverer,
	bugfixes.LowerCaseHeaders,
)
handler := system.Handler(mux)
```

### CORS

Enable CORS with domain-specific configuration:

```go
system := bugfixes.NewMiddleware()
system.AddAllowedOrigins("https://example.com", "https://app.example.com")
system.AddAllowedMethods("GET", "POST", "PUT", "DELETE")
system.AddAllowedHeaders("Authorization", "X-Custom-Header")
system.AddMiddleware(system.CORS)
handler := system.Handler(mux)
```

Multiple domains are supported. Requests from unlisted origins return `403 Forbidden`.

### Security Headers

Enable security headers with production-safe defaults:

```go
system := bugfixes.NewMiddleware()
system.SetSecure(true)
system.AddMiddleware(system.Secure)
handler := system.Handler(mux)
```

This sets the following headers by default:

| Header | Value |
|--------|-------|
| X-Frame-Options | DENY |
| X-Content-Type-Options | nosniff |
| X-XSS-Protection | 1; mode=block |
| Strict-Transport-Security | max-age=31536000; includeSubDomains |
| Content-Security-Policy | default-src 'self' |
| Referrer-Policy | strict-origin-when-cross-origin |

#### Domain-Specific Configuration

Override defaults per-domain:

```go
import "time"

system := bugfixes.NewMiddleware()
system.SetSecure(true)

// Allow iframes on main site
system.AddSecureConfig("example.com", bugfixes.SecureConfig{
	XFrameOptions: "SAMEORIGIN",
	CSP:           "default-src 'self'; script-src 'self' 'unsafe-inline'",
})

// Stricter CSP for API subdomain
system.AddSecureConfig("api.example.com", bugfixes.SecureConfig{
	XFrameOptions: "DENY",
	CSP:           "default-src 'self'",
	HSTSEnabled:   true,
	HSTSMaxAge:    730 * 24 * time.Hour, // 2 years
})

system.AddMiddleware(system.Secure)
handler := system.Handler(mux)
```

#### HSTS Options

Configure HSTS with additional options:

```go
import "time"

system.AddSecureConfig("example.com", bugfixes.SecureConfig{
	HSTSEnabled:            true,
	HSTSMaxAge:            180 * 24 * time.Hour, // 6 months
	HSTSIncludeSubdomains: true,
	HSTSPreload:           true,
})
```

#### Disabling Headers

Set empty string or `false` to disable individual headers:

```go
system.AddSecureConfig("example.com", bugfixes.SecureConfig{
	XFrameOptions:        "",    // removes X-Frame-Options
	XContentTypeOptions:  false, // removes X-Content-Type-Options
	XXSSProtection:       "",    // removes X-XSS-Protection
})
```

## Development

Common repo tasks:

```bash
just fmt
just lint
just test
just test-race
just check
```
