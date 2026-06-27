package httperror

import "net/http"

func SignInInvalidCredentialsError() *HttpError {
	// Keep message intentionally vague for security
	return &HttpError{Status: http.StatusUnauthorized, Message: "Invalid email or password.", Code: "INTERNAL_SERVER_ERROR"}
}
