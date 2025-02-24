package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestJSONMiddleware tests the JSONMiddleware function
func TestJSONMiddleware(t *testing.T) {
	// Create a test handler that will check the Content-Type header
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the Content-Type is set to application/json
		if w.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
		}
	})

	// Create a request to pass to our handler
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create the middleware and call it
	middleware := JSONMiddleware(testHandler)
	middleware.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", status)
	}

	// Check if the Content-Type header is set correctly
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
}
