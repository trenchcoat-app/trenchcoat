package api_error

import "net/http"

func SignOutSessionNotFoundError() *ApiError {
	return &ApiError{
		Status:  http.StatusUnauthorized,
		Message: "Session not found or already expired.",
		Code:    "UNAUTHORIZED",
	}
}

func SignOutUnauthorizedError(msg string) *ApiError {
	return &ApiError{
		Status:  http.StatusUnauthorized,
		Message: msg,
		Code:    "UNAUTHORIZED",
	}
}
