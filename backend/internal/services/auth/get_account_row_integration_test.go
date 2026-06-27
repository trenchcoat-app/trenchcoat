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

func TestGetAccountRow_Success(t *testing.T) {
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
		Password: "correct-password",
	}

	account, httpErr := svc.GetAccountRow(c, body)
	require.Nil(t, httpErr)
	require.NotNil(t, account)
	assert.Equal(t, "Test User", account.DisplayName)
	assert.Equal(t, "active", account.Status)
}

func TestGetAccountRow_InvalidCredentials(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)

	t.Run("wrong password", func(t *testing.T) {
		t.Parallel()
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
		_, httpErr := svc.GetAccountRow(c, body)
		require.NotNil(t, httpErr)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Status)
	})

	t.Run("unknown email", func(t *testing.T) {
		t.Parallel()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

		body := api.SignInJSONRequestBody{
			Email:    testutil.NewEmail(),
			Password: "any-password",
		}
		_, httpErr := svc.GetAccountRow(c, body)
		require.NotNil(t, httpErr)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Status)
	})
}

func TestGetAccountRow_DisabledAccount(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)
	email := testutil.NewEmail()
	hash := testutil.HashPassword(t, "any-password")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	var id uuid.UUID
	err := pool.QueryRow(
		context.Background(),
		`INSERT INTO account (email, display_name, password_hash, status) VALUES ($1, $2, $3, 'disabled') RETURNING id`,
		string(email), "Disabled User", hash,
	).Scan(&id)
	require.NoError(t, err)

	body := api.SignInJSONRequestBody{
		Email:    email,
		Password: "any-password",
	}

	_, httpErr := svc.GetAccountRow(c, body)
	require.NotNil(t, httpErr)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Status)
}
