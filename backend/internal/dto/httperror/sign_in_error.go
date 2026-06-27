package httperror

import "net/http"

func SignInInvalidCredentialsError() *HTTPError {
	// Keep message intentionally vague for security
	return &HTTPError{Status: http.StatusUnauthorized, Message: "Invalid email or password.", Code: "INVALID_CREDENTIALS"}
}
