package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"trenchcoat/internal/services/auth"
	"trenchcoat/internal/services/testutil"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
)

func TestSignOut_Success(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)
	email := testutil.NewEmail()
	accountID := testutil.SeedAccount(t, pool, string(email), "Signout User", "password")
	token := testutil.SeedSession(t, pool, accountID)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	httpErr := svc.SignOut(c, token)
	require.Nil(t, httpErr)

	var count int
	err := pool.QueryRow(
		context.Background(), `SELECT COUNT(*) FROM session WHERE token = $1`, token,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestSignOut_NotFound(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	token := uuid.New()
	httpErr := svc.SignOut(c, token)
	require.NotNil(t, httpErr)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Status)
	assert.Equal(t, "UNAUTHORIZED", httpErr.Code)
}
