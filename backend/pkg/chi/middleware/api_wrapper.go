package middleware

import (
	"backend/pkg/chi/response"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
)

// APIWrapper is a middleware function that wraps an HTTP handler and provides error handling and response building.
// It takes a handler function that returns a response.Body and returns an http.HandlerFunc.
func APIWrapper(handler func(w http.ResponseWriter, r *http.Request) response.Body) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := zerolog.Ctx(r.Context())
		data := func() response.Body {
			defer func() {
				if e := recover(); e != nil {
					err, ok := e.(error)
					if !ok {
						err = fmt.Errorf("%v", e)
					}
					l.Error().Err(err).Msgf("recovered from panic: %s", debug.Stack())
				}
			}()
			return handler(w, r)
		}()

		Build(data).Render(w, r)
	}
}

// Build constructs a response.Body based on the provided data.
// It handles different types of responses and returns the appropriate response body.
func Build(data any) response.Body {
	switch d := data.(type) {
	case response.HTTPSuccess:
		return d
	case response.HTTPError:
		return d
	default:
		return response.InternalServerError("")
	}
}
