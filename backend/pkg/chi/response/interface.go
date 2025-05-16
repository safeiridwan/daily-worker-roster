package response

import "net/http"

// Body represents an interface for rendering response bodies.
type Body interface {
	// Render writes the response body to the provided http.ResponseWriter.
	// It takes the http.ResponseWriter and *http.Request as parameters and returns an error if any occurs during rendering.
	Render(w http.ResponseWriter, r *http.Request) error
}
