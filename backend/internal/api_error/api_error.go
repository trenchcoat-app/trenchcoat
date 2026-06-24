package api_error

import (
	"net/http"
	"trenchcoat/internal/api"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

func (e *ApiError) Error() string {
	return e.Message
}

// Write the matching JSON response of the ApiError to the Gin context
func HandleApiError(c *gin.Context, err ApiError) {
	c.JSON(err.Status, api.ErrorResponse{
		Code:    err.Code,
		Message: err.Message,
	})
}

func InternalServerError(msg string) *ApiError {
	return &ApiError{Status: http.StatusInternalServerError, Message: msg, Code: "INTERNAL_SERVER_ERROR"}
}
