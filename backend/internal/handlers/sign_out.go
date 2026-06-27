package handlers

import (
	"net/http"
	"trenchcoat/internal/httperror"
	"trenchcoat/internal/services/cookie"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) SignOut(c *gin.Context) {
	var tokenUUID uuid.UUID
	var apiErr *httperror.HttpError

	tokenStr, err := cookie.GetSessionToken(c)
	if err == nil && tokenStr != "" {
		tokenUUID, err = uuid.Parse(tokenStr)
		if err != nil {
			apiErr = httperror.SignOutUnauthorizedError("Invalid session token.")
		}
	} else {
		// Authorization header fallback
		authHeader := c.GetHeader("Authorization")
		tokenUUID, apiErr = s.AuthService.ParseAuthToken(authHeader)
	}

	if apiErr != nil {
		httperror.HandleHttpError(c, *apiErr)
		return
	}

	apiErr = s.AuthService.SignOut(c, tokenUUID)
	if apiErr != nil {
		httperror.HandleHttpError(c, *apiErr)
		return
	}

	cookie.ClearSessionCookie(c)

	c.Status(http.StatusNoContent)
}
