package api_error

import "net/http"

func SignUpEmailAlreadyExistsError() *ApiError {
	return &ApiError{
		Status:  http.StatusConflict,
		Message: "An account with this email address already exists.",
		Code:    "EMAIL_ALREADY_EXISTS",
	}
}
