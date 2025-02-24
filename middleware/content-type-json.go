package middleware

import "net/http"

// JSONMiddleware sets the Content-Type header to application/json for all responses
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure Content-Type is set for ALL responses (errors or success)
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
