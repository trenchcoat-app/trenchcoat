package api_error

import "net/http"

func SignInInvalidCredentialsError() *ApiError {
	// Keep message intentionally vague for security
	return &ApiError{Status: http.StatusUnauthorized, Message: "Invalid email or password.", Code: "INTERNAL_SERVER_ERROR"}
}
