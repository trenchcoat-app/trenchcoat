package auth

import (
	"time"
	"trenchcoat/config"
	"trenchcoat/internal/api_error"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type Session struct {
	SessionToken *openapi_types.UUID
	ExpiresAt    *time.Time
}

func (auth *AuthService) CreateSession(c *gin.Context, account AccountRow) (session Session, apiErr *api_error.ApiError) {
	session.ExpiresAt = auth.GetNewSessionExpireTime(config.AppConfig.SESSION_EXPIRY_SECONDS)

	sql := `
		INSERT INTO session (expires_at, ip_address, user_agent, account_id)
		VALUES ($1, $2, $3, $4)
		RETURNING token, expires_at
	`

	err := auth.DB.QueryRow(
		c.Request.Context(),
		sql,
		session.ExpiresAt,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		account.ID,
	).Scan(&session.SessionToken, &session.ExpiresAt)

	if err != nil {
		apiErr = api_error.InternalServerError("Failed to create session: " + err.Error())
	}

	return
}

func (auth *AuthService) GetNewSessionExpireTime(offsetSeconds int) *time.Time {
	expireTime := time.Now().Add(time.Duration(offsetSeconds * int(time.Second)))
	return &expireTime
}
