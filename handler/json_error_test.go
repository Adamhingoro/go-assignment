package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONError(t *testing.T) {
	// Create a response recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the JSONError function with test parameters
	JSONError(recorder, http.StatusBadRequest, "Test error message")

	// Check the status code
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, status)
	}

	// Check the response body
	expected := map[string]string{"error": "Test error message"}
	expectedBody, _ := json.Marshal(expected) // Marshal the expected response

	// Normalize both expected and actual response bodies
	var expectedCompact, actualCompact bytes.Buffer
	json.Compact(&expectedCompact, expectedBody)
	json.Compact(&actualCompact, recorder.Body.Bytes())

	if expectedCompact.String() != actualCompact.String() {
		t.Errorf("expected body %s, got %s", expectedCompact.String(), actualCompact.String())
	}
}
