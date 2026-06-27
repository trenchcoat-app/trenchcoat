package cookie_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"trenchcoat/internal/services/cookie"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
)

func TestSetSessionCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	token := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour)

	cookie.SetSessionCookie(c, token, expiresAt)

	result := w.Result()
	cookies := result.Cookies()
	require.Len(t, cookies, 1)

	sessionCookie := cookies[0]
	assert.Equal(t, "sid", sessionCookie.Name)
	assert.Equal(t, token.String(), sessionCookie.Value)
	assert.Equal(t, "/", sessionCookie.Path)
	assert.True(t, sessionCookie.HttpOnly)
	assert.False(t, sessionCookie.Secure)
	assert.Equal(t, "", sessionCookie.Domain)
	assert.Greater(t, sessionCookie.MaxAge, 0)
}

func TestClearSessionCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	cookie.ClearSessionCookie(c)

	result := w.Result()
	cookies := result.Cookies()
	require.Len(t, cookies, 1)

	sessionCookie := cookies[0]
	assert.Equal(t, "sid", sessionCookie.Name)
	assert.Equal(t, "", sessionCookie.Value)
	assert.Equal(t, "/", sessionCookie.Path)
	assert.True(t, sessionCookie.HttpOnly)
	assert.False(t, sessionCookie.Secure)
	assert.Equal(t, "", sessionCookie.Domain)
	assert.Less(t, sessionCookie.MaxAge, 0)
	assert.True(t, sessionCookie.Expires.Before(time.Now()))
}

func TestGetSessionToken_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "sid",
		Value: "test-token-value",
	})
	c.Request = req

	token, err := cookie.GetSessionToken(c)
	require.NoError(t, err)
	assert.Equal(t, "test-token-value", token)
}

func TestGetSessionToken_Missing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("GET", "/", nil)
	c.Request = req

	_, err := cookie.GetSessionToken(c)
	require.Error(t, err)
	assert.ErrorIs(t, err, http.ErrNoCookie)
}
