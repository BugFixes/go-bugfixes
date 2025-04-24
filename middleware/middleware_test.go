package middleware_test

import (
	"github.com/bugfixes/go-bugfixes/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBugFixes(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Error(err)
		}
	})

	// Create a request to pass to our handler
	// We don't have any query parameters so we'll pass 'nil' as the third parameter
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Apply middleware to the handler
	middlewareHandler := middleware.DefaultMiddleware(handler)

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `OK`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v, want %v", rr.Body.String(), expected)
	}
}

func TestBugFixesLowerCase(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Error(err)
		}
	})

	// Create a request to pass to our handler
	// We don't have any query parameters so we'll pass 'nil' as the third parameter
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Bugfixes-Version", "1.0.0")

	// Apply middleware to the handler
	middlewareHandler := middleware.DefaultMiddleware(handler)

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(rr, req)

	// make sure the lowercase works
	if req.Header.Get("X-Bugfixes-Version") != "1.0.0" || req.Header.Get("x-bugfixes-version") != "1.0.0" {
		t.Fatal("handler returned wrong header")
	}
}
