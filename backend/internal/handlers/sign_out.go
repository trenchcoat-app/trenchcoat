package handlers

import (
	"net/http"
	"trenchcoat/internal/dto/httperror"
	"trenchcoat/internal/services/cookie"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) SignOut(c *gin.Context) {
	var tokenUUID uuid.UUID
	var httpErr *httperror.HttpError

	tokenStr, err := cookie.GetSessionToken(c)
	if err == nil && tokenStr != "" {
		tokenUUID, err = uuid.Parse(tokenStr)
		if err != nil {
			httpErr = httperror.SignOutUnauthorizedError("Invalid session token.")
		}
	} else {
		// Authorization header fallback
		authHeader := c.GetHeader("Authorization")
		tokenUUID, httpErr = s.AuthService.ParseAuthToken(authHeader)
	}

	if httpErr != nil {
		httperror.HandleHttpError(c, httpErr)
		return
	}

	httpErr = s.AuthService.SignOut(c, tokenUUID)
	if httpErr != nil {
		httperror.HandleHttpError(c, httpErr)
		return
	}

	cookie.ClearSessionCookie(c)

	c.Status(http.StatusNoContent)
}
