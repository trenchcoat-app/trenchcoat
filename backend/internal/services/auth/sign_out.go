package auth

import (
	"strings"
	"trenchcoat/internal/httperror"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (auth *AuthService) ParseAuthToken(authHeader string) (uuid.UUID, *httperror.HttpError) {
	if authHeader == "" {
		return uuid.Nil, httperror.SignOutUnauthorizedError("Missing authorization token.")
	}

	tokenStr := authHeader
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		tokenStr = authHeader[7:]
	}
	tokenStr = strings.TrimSpace(tokenStr)

	if tokenStr == "" {
		return uuid.Nil, httperror.SignOutUnauthorizedError("Missing authorization token.")
	}

	tokenUUID, err := uuid.Parse(tokenStr)
	if err != nil {
		return uuid.Nil, httperror.SignOutUnauthorizedError("Invalid authorization token format.")
	}

	return tokenUUID, nil
}

func (auth *AuthService) SignOut(c *gin.Context, tokenUUID uuid.UUID) *httperror.HttpError {
	sql := `
		DELETE FROM session
		WHERE token = $1
	`

	cmdTag, err := auth.DB.Exec(c.Request.Context(), sql, tokenUUID)
	if err != nil {
		return httperror.InternalServerError("Failed to invalidate session: " + err.Error())
	}

	if cmdTag.RowsAffected() == 0 {
		return httperror.SignOutSessionNotFoundError()
	}

	return nil
}
