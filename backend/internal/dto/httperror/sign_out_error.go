package httperror

import "net/http"

func SignOutSessionNotFoundError() *HTTPError {
	return &HTTPError{
		Status:  http.StatusUnauthorized,
		Message: "Session not found or already expired.",
		Code:    "UNAUTHORIZED",
	}
}

func SignOutUnauthorizedError(msg string) *HTTPError {
	return &HTTPError{
		Status:  http.StatusUnauthorized,
		Message: msg,
		Code:    "UNAUTHORIZED",
	}
}
