package handlers

import (
	"net/http"
	"strings"
	"trenchcoat/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) SignOut(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Missing authorization token.",
		})
		return
	}

	tokenStr := authHeader
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		tokenStr = authHeader[7:]
	}
	tokenStr = strings.TrimSpace(tokenStr)

	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Missing authorization token.",
		})
		return
	}

	tokenUUID, err := uuid.Parse(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Invalid authorization token format.",
		})
		return
	}

	cmdTag, err := s.DB.Exec(c.Request.Context(), "DELETE FROM session WHERE token = $1", tokenUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to invalidate session: " + err.Error(),
		})
		return
	}

	if cmdTag.RowsAffected() == 0 {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Session not found or already expired.",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
