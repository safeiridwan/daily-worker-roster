package response

import (
	"encoding/json"
	"net/http"
)

// HTTPSuccess represents a successful HTTP response.
type HTTPSuccess struct {
	HTTPStatus int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

// Render writes the HTTP response.
func (e HTTPSuccess) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.HTTPStatus)
	return json.NewEncoder(w).Encode(e)
}

// NewHTTPSuccess creates a new HTTPSuccess instance.
func NewHTTPSuccess(httpStatus int, message string, data ...interface{}) HTTPSuccess {
	httpSuccess := HTTPSuccess{
		HTTPStatus: httpStatus,
		Message:    message,
	}
	if len(data) > 0 {
		httpSuccess.Data = data[0]
	}
	return httpSuccess
}

// OK creates a new HTTPSuccess instance with HTTP status 200 OK.
func OK(msg string, data ...interface{}) HTTPSuccess {
	resp := HTTPSuccess{
		HTTPStatus: http.StatusOK,
		Message:    msg,
	}
	if len(data) > 0 {
		resp.Data = data[0]
	}
	return resp
}
