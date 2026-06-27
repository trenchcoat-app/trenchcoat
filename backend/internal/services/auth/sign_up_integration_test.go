package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"trenchcoat/internal/api"
	"trenchcoat/internal/services/auth"
	"trenchcoat/internal/services/testutil"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestCreateAccount_Success(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)
	email := testutil.NewEmail()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	hash := testutil.HashPassword(t, "secure-password")
	id, httpErr := svc.CreateAccount(c, string(email), "New User", hash)
	require.Nil(t, httpErr)

	var storedEmail string
	var displayName string
	err := pool.QueryRow(
		context.Background(),
		`SELECT email, display_name FROM account WHERE id = $1`, id,
	).Scan(&storedEmail, &displayName)
	require.NoError(t, err)
	assert.Equal(t, string(email), storedEmail)
	assert.Equal(t, "New User", displayName)
}

func TestCreateAccount_DuplicateEmail(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)
	email := testutil.NewEmail()
	testutil.SeedAccount(t, pool, string(email), "First User", "password1")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	hash := testutil.HashPassword(t, "password2")
	_, httpErr := svc.CreateAccount(c, string(email), "Second User", hash)
	require.NotNil(t, httpErr)
	assert.Equal(t, http.StatusInternalServerError, httpErr.Status)
}

func TestSignUp_Success(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)
	email := testutil.NewEmail()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	autoSignIn := false
	body := api.SignUpJSONRequestBody{
		Email:       email,
		Password:    "secure-password",
		DisplayName: "New User",
		AutoSignIn:  &autoSignIn,
	}

	resp, httpErr := svc.SignUp(c, body)
	require.Nil(t, httpErr)
	require.NotNil(t, resp)
	assert.Equal(t, "New User", *resp.Account.DisplayName)
	assert.Equal(t, email, resp.Account.Email)
}

func TestSignUp_DuplicateEmail(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)
	email := testutil.NewEmail()
	testutil.SeedAccount(t, pool, string(email), "First User", "password1")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := api.SignUpJSONRequestBody{
		Email:       email,
		Password:    "another-password",
		DisplayName: "Second User",
	}

	_, httpErr := svc.SignUp(c, body)
	require.NotNil(t, httpErr)
	assert.Equal(t, http.StatusConflict, httpErr.Status)
	assert.Equal(t, "EMAIL_ALREADY_EXISTS", httpErr.Code)
}

func TestSignUp_AutoSignIn(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)
	email := testutil.NewEmail()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	autoSignIn := true
	body := api.SignUpJSONRequestBody{
		Email:       email,
		Password:    "secure-password",
		DisplayName: "Auto User",
		AutoSignIn:  &autoSignIn,
	}

	resp, httpErr := svc.SignUp(c, body)
	require.Nil(t, httpErr)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Session)
	require.NotNil(t, resp.Session.SessionToken)
}
