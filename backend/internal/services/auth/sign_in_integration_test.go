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
	"github.com/google/uuid"
)

func TestSignIn_Success(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)
	email := testutil.NewEmail()
	password := "the-password"
	testutil.SeedAccount(t, pool, string(email), "Signin User", password)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := api.SignInJSONRequestBody{
		Email:    email,
		Password: password,
	}

	resp, httpErr := svc.SignIn(c, body)
	require.Nil(t, httpErr)
	require.NotNil(t, resp)
	assert.Equal(t, "Signin User", *resp.Account.DisplayName)
	assert.Equal(t, email, resp.Account.Email)
	require.NotNil(t, resp.Session.SessionToken)
}

func TestSignIn_InvalidCredentials(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)
	email := testutil.NewEmail()
	testutil.SeedAccount(t, pool, string(email), "Test User", "correct-password")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := api.SignInJSONRequestBody{
		Email:    email,
		Password: "wrong-password",
	}

	_, httpErr := svc.SignIn(c, body)
	require.NotNil(t, httpErr)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Status)
}

func TestSignIn_DisabledAccount(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)
	email := testutil.NewEmail()
	hash := testutil.HashPassword(t, "password")

	var id uuid.UUID
	err := pool.QueryRow(
		context.Background(),
		`INSERT INTO account (email, display_name, password_hash, status) VALUES ($1, $2, $3, 'disabled') RETURNING id`,
		string(email), "Disabled User", hash,
	).Scan(&id)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := api.SignInJSONRequestBody{
		Email:    email,
		Password: "password",
	}

	_, httpErr := svc.SignIn(c, body)
	require.NotNil(t, httpErr)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Status)
}
