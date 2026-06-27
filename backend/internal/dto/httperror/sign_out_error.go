package httperror

import "net/http"

func SignOutSessionNotFoundError() *HttpError {
	return &HttpError{
		Status:  http.StatusUnauthorized,
		Message: "Session not found or already expired.",
		Code:    "UNAUTHORIZED",
	}
}

func SignOutUnauthorizedError(msg string) *HttpError {
	return &HttpError{
		Status:  http.StatusUnauthorized,
		Message: msg,
		Code:    "UNAUTHORIZED",
	}
}
