package cookie

import (
	"net/http"
	"time"
	"trenchcoat/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetSessionCookie(c *gin.Context, token uuid.UUID, expiresAt time.Time) {
	c.SetCookie(
		config.AppConfig.COOKIE_NAME,
		token.String(),
		int(time.Until(expiresAt).Seconds()),
		config.AppConfig.COOKIE_PATH,
		config.AppConfig.COOKIE_DOMAIN,
		config.AppConfig.COOKIE_SECURE,
		true,
	)
}

func ClearSessionCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.AppConfig.COOKIE_NAME,
		Value:    "",
		Path:     config.AppConfig.COOKIE_PATH,
		Domain:   config.AppConfig.COOKIE_DOMAIN,
		Secure:   config.AppConfig.COOKIE_SECURE,
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

func GetSessionToken(c *gin.Context) (string, error) {
	return c.Cookie(config.AppConfig.COOKIE_NAME)
}
