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

## Development

Common repo tasks:

```bash
just fmt
just lint
just test
just test-race
just check
```
