package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-playground/validator/v10"
)

// HTTPError represents an HTTP error response.
type HTTPError struct {
	HTTPStatus int         `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
}

// Render writes the error response to the http.ResponseWriter.
func (e HTTPError) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.HTTPStatus)
	return json.NewEncoder(w).Encode(e)
}

// NewHTTPError creates a new HTTPError with the given status, message, and optional details.
func NewHTTPError(httpStatus int, msg string, details ...interface{}) HTTPError {
	resp := HTTPError{
		HTTPStatus: httpStatus,
		Message:    msg,
	}
	if len(details) > 0 {
		resp.Details = details[0]
	}
	return resp
}

// defaultMessages contains default error messages for common HTTP status codes.
var defaultMessages = map[int]string{
	http.StatusInternalServerError: "We encountered an error while processing your request.",
	http.StatusNotFound:            "The requested resource was not found.",
	http.StatusUnauthorized:        "You are not authenticated to perform the requested action.",
	http.StatusForbidden:           "You are not authorized to perform the requested action.",
	http.StatusBadRequest:          "Your request is in a bad format.",
	http.StatusMethodNotAllowed:    "Your request method is not supported by the server.",
}

// InternalServerError creates an HTTPError for internal server errors.
func InternalServerError(msg string) HTTPError {
	return NewHTTPError(http.StatusInternalServerError, getMessage(msg, http.StatusInternalServerError))
}

// NotFound creates an HTTPError for not found resources.
func NotFound(msg string) HTTPError {
	return NewHTTPError(http.StatusNotFound, getMessage(msg, http.StatusNotFound))
}

// Unauthorized creates an HTTPError for unauthorized requests.
func Unauthorized(msg string) HTTPError {
	return NewHTTPError(http.StatusUnauthorized, getMessage(msg, http.StatusUnauthorized))
}

// Forbidden creates an HTTPError for forbidden requests.
func Forbidden(msg string) HTTPError {
	return NewHTTPError(http.StatusForbidden, getMessage(msg, http.StatusForbidden))
}

// BadRequest creates an HTTPError for bad requests.
func BadRequest(msg string) HTTPError {
	return NewHTTPError(http.StatusBadRequest, getMessage(msg, http.StatusBadRequest))
}

// MethodNotAllowed creates an HTTPError for method not allowed requests.
func MethodNotAllowed(msg string) HTTPError {
	return NewHTTPError(http.StatusMethodNotAllowed, getMessage(msg, http.StatusMethodNotAllowed))
}

// getMessage returns the provided message if not empty, otherwise returns the default message for the given status.
func getMessage(msg string, status int) string {
	if msg == "" {
		return defaultMessages[status]
	}
	return msg
}

// invalidField represents an invalid field in a request.
type invalidField struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// InvalidInput creates an HTTPError for invalid input with details about the invalid fields.
func InvalidInput(err error) HTTPError {
	var details []invalidField

	switch e := err.(type) {
	case validation.ValidationErrors:
		for _, field := range e {
			details = append(details, invalidField{
				Field: field.Field(),
				Error: getInvalidInputMsg(field),
			})
		}
	}

	return NewHTTPError(
		http.StatusBadRequest,
		"There is some problem with the data you submitted.",
		details,
	)
}

// getInvalidInputMsg returns a user-friendly error message for invalid input fields.
func getInvalidInputMsg(fieldErr validation.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "This field is required"
	case "gte":
		return fmt.Sprintf("This field must be greater than or equal to %s", fieldErr.Param())
	case "lte":
		return fmt.Sprintf("This field must be less than or equal to %s", fieldErr.Param())
	case "datetime":
		return fmt.Sprintf("This field must be in the format %s", fieldErr.Param())
	default:
		fmt.Println(fieldErr.Tag())
		return fieldErr.Error()
	}
}
