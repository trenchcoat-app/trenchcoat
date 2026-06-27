package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"trenchcoat/internal/api"
	"trenchcoat/internal/dto/httperror"

	"trenchcoat/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

var signOutRoute = "/api/v1/auth/sign-out"

func TestSignOut_SuccessViaCookie(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	token := uuid.New()

	mockAuth := NewMockAuthServiceInterface(t)
	mockAuth.EXPECT().
		SignOut(mock.Anything, token).
		Return(nil)

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("POST", signOutRoute, nil)
	req.AddCookie(&http.Cookie{
		Name:  "sid",
		Value: token.String(),
	})
	c.Request = req

	srv.SignOut(c)

	assert.Equal(t, http.StatusNoContent, c.Writer.Status())

	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, "sid", cookies[0].Name)
	assert.Equal(t, "", cookies[0].Value)
	assert.Less(t, cookies[0].MaxAge, 0)
	assert.True(t, cookies[0].Expires.Before(time.Now()))
}

func TestSignOut_SuccessViaAuthHeader(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	token := uuid.New()

	mockAuth := NewMockAuthServiceInterface(t)
	mockAuth.EXPECT().
		ParseAuthToken("Bearer "+token.String()).
		Return(token, nil)
	mockAuth.EXPECT().
		SignOut(mock.Anything, token).
		Return(nil)

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("POST", signOutRoute, nil)
	c.Request = req
	c.Request.Header.Set("Authorization", "Bearer "+token.String())

	srv.SignOut(c)

	assert.Equal(t, http.StatusNoContent, c.Writer.Status())

	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, "sid", cookies[0].Name)
	assert.Equal(t, "", cookies[0].Value)
	assert.Less(t, cookies[0].MaxAge, 0)
}

func TestSignOut_MissingAuth(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	mockAuth := NewMockAuthServiceInterface(t)
	mockAuth.EXPECT().
		ParseAuthToken("").
		Return(uuid.Nil, httperror.SignOutUnauthorizedError("Missing authorization token."))

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", signOutRoute, nil)

	srv.SignOut(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "UNAUTHORIZED", resp.Code)
	assert.Equal(t, "Missing authorization token.", resp.Message)
}

func TestSignOut_InvalidCookieUUID(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	mockAuth := NewMockAuthServiceInterface(t)

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("POST", signOutRoute, nil)
	req.AddCookie(&http.Cookie{
		Name:  "sid",
		Value: "not-a-uuid",
	})
	c.Request = req

	srv.SignOut(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "UNAUTHORIZED", resp.Code)
	assert.Equal(t, "Invalid session token.", resp.Message)
}

func TestSignOut_SessionNotFound(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	token := uuid.New()

	mockAuth := NewMockAuthServiceInterface(t)
	mockAuth.EXPECT().
		SignOut(mock.Anything, token).
		Return(httperror.SignOutSessionNotFoundError())

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("POST", signOutRoute, nil)
	req.AddCookie(&http.Cookie{
		Name:  "sid",
		Value: token.String(),
	})
	c.Request = req

	srv.SignOut(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "UNAUTHORIZED", resp.Code)
	assert.Equal(t, "Session not found or already expired.", resp.Message)
}
