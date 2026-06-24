package handlers

import (
	"net/http"
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

	apiErr = s.AuthService.SignOut(c, tokenUUID)
	if apiErr != nil {
		api_error.HandleApiError(c, *apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}
