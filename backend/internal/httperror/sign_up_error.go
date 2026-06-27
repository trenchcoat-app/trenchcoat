package httperror

import "net/http"

func SignUpEmailAlreadyExistsError() *HttpError {
	return &HttpError{
		Status:  http.StatusConflict,
		Message: "An account with this email address already exists.",
		Code:    "EMAIL_ALREADY_EXISTS",
	}
}
