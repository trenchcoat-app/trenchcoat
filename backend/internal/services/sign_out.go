package services

import (
	"strings"
	"trenchcoat/internal/api_error"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (auth *AuthService) ParseAuthToken(authHeader string) (uuid.UUID, *api_error.ApiError) {
	if authHeader == "" {
		return uuid.Nil, api_error.SignOutUnauthorizedError("Missing authorization token.")
	}

	tokenStr := authHeader
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		tokenStr = authHeader[7:]
	}
	tokenStr = strings.TrimSpace(tokenStr)

	if tokenStr == "" {
		return uuid.Nil, api_error.SignOutUnauthorizedError("Missing authorization token.")
	}

	tokenUUID, err := uuid.Parse(tokenStr)
	if err != nil {
		return uuid.Nil, api_error.SignOutUnauthorizedError("Invalid authorization token format.")
	}

	return tokenUUID, nil
}

func (auth *AuthService) SignOut(c *gin.Context, tokenUUID uuid.UUID) *api_error.ApiError {
	sql := `
		DELETE FROM session
		WHERE token = $1
	`

	cmdTag, err := auth.DB.Exec(c.Request.Context(), sql, tokenUUID)
	if err != nil {
		return api_error.InternalServerError("Failed to invalidate session: " + err.Error())
	}

	if cmdTag.RowsAffected() == 0 {
		return api_error.SignOutSessionNotFoundError()
	}

	return nil
}
