package httperror

import (
	"net/http"
	"trenchcoat/internal/api"

	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

// Write the matching JSON response of the HttpError to the Gin context
func HandleHttpError(c *gin.Context, err *HTTPError) {
	c.JSON(err.Status, api.ErrorResponse{
		Code:    err.Code,
		Message: err.Message,
	})
}

func InternalServerError(msg string) *HTTPError {
	return &HTTPError{Status: http.StatusInternalServerError, Message: msg, Code: "INTERNAL_SERVER_ERROR"}
}
