package middleware_test

import (
	"context"
	"github.com/bugfixes/go-bugfixes/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	// Define the next HTTP handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if reqID := middleware.GetReqID(r.Context()); reqID == "" {
			t.Error("Request ID missing from context")
		}
	})

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Apply middleware
	handler := middleware.RequestID(nextHandler)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
}

func TestGetReqID(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.RequestIDKey, "12345")
	reqID := middleware.GetReqID(ctx)

	if reqID != "12345" {
		t.Errorf("Incorrect request ID in context, got: %s, want: %s.", reqID, "12345")
	}
}

func TestNextRequestID(t *testing.T) {
	id1 := middleware.NextRequestID()
	id2 := middleware.NextRequestID()

	if id2 != id1+1 {
		t.Errorf("Incorrect sequential IDs, got: %d and %d, want: %d and %d.", id1, id2, id1, id1+1)
	}
}
