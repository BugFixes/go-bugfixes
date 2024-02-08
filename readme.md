# Bugfixes Agent for Go

There are a number of tools in this package
  - Logging
  - Middleware

Logging allows for the following log levels
  - Log / Logf / Printf / Sprintf
  - Info / Infof
  - Debug / Debugf
  - Error / Errorf

### Requirements
This is the go agent for bugfixes you will need the following in your environment

- BUGFIXES_AGENT_KEY
- BUGFIXES_AGENT_SECRET

#### If you wish to keep the crash local (e.g. for testing purposes)
- BUGFIXES_LOCAL_ONLY

___
## Middleware
If you use Chi, Gorilla or multiple other http server packages that support middleware, you can use the middleware
```go
package main
import (
	"net/http"
	bugfixes "github.com/bugfixes/go-bugfixes/middleware"
	buglog "github.com/bugfixes/go-bugfixes/logs"
  "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
  r := chi.NewRouter()
  r.Use(middleware.Timeout(60 * time.Second))
  r.Use(bugfixes.Middlware)

  r.Route("/", tester)
  if err := http.ListenAndServe(":8080", r); err != nil {
  	return buglog.Errorf("failed to start listener: %v", err)
  }

  return nil
}

func tester(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
```
