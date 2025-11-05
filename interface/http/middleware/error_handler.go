package middleware

import (
	"fmt"
	"net/http"
)

// HTTPError represents a standard HTTP error response
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewHTTPError creates a new HTTPError
func NewHTTPError(code int, message string, details ...string) *HTTPError {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}

	return &HTTPError{
		Code:    code,
		Message: message,
		Details: detail,
	}
}

// ErrorHandler provides error handling utilities
type ErrorHandler struct{}

// NewErrorHandler creates a new ErrorHandler
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleError converts an error to an HTTP error response
func (h *ErrorHandler) HandleError(err error) *HTTPError {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Map specific error messages to HTTP status codes
	switch {
	case errMsg == "task not found" || errMsg == "project not found" || errMsg == "user not found":
		return NewHTTPError(http.StatusNotFound, "Resource not found", errMsg)

	case errMsg == "invalid task id" || errMsg == "invalid project id" || errMsg == "invalid user id":
		return NewHTTPError(http.StatusBadRequest, "Invalid input", errMsg)

	case errMsg == "task is already assigned":
		return NewHTTPError(http.StatusConflict, "Task already assigned", errMsg)

	case errMsg == "cannot transition" || errMsg == "invalid status transition":
		return NewHTTPError(http.StatusBadRequest, "Invalid state transition", errMsg)

	default:
		return NewHTTPError(
			http.StatusInternalServerError,
			"Internal server error",
			fmt.Sprintf("An unexpected error occurred: %v", err),
		)
	}
}