package cookie

import (
	"time"
	"trenchcoat/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetSessionCookie(c *gin.Context, token uuid.UUID, expiresAt time.Time) {
	c.SetCookie(
		config.AppConfig.SESSION_COOKIE_NAME,
		token.String(),
		int(time.Until(expiresAt).Seconds()),
		config.AppConfig.SESSION_COOKIE_PATH,
		config.AppConfig.SESSION_COOKIE_DOMAIN,
		config.AppConfig.SESSION_COOKIE_SECURE,
		true,
	)
}

func ClearSessionCookie(c *gin.Context) {
	c.SetCookie(
		config.AppConfig.SESSION_COOKIE_NAME,
		"",
		-1,
		config.AppConfig.SESSION_COOKIE_PATH,
		config.AppConfig.SESSION_COOKIE_DOMAIN,
		config.AppConfig.SESSION_COOKIE_SECURE,
		true,
	)
}

func GetSessionToken(c *gin.Context) (string, error) {
	return c.Cookie(config.AppConfig.SESSION_COOKIE_NAME)
}
