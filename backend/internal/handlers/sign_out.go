package handlers

import (
	"errors"
	"net/http"
	"trenchcoat/internal/api"
	"trenchcoat/internal/api_error"

	"github.com/gin-gonic/gin"
)

func (s *Server) SignOut(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenUUID, apiErr := s.AuthService.ParseAuthToken(authHeader)

	if apiErr != nil {
		api_error.HandleApiError(c, *apiErr)
		return
	}

	err := s.AuthService.SignOut(c, tokenUUID)

	if err != nil {
		var apiErr *api_error.ApiError
		if errors.As(err, &apiErr) {
			api_error.HandleApiError(c, *apiErr)
			return
		}

		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to invalidate session: " + err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
