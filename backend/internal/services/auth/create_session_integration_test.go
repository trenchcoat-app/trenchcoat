package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"trenchcoat/internal/services/auth"
	"trenchcoat/internal/services/testutil"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestCreateSession_Success(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	svc := auth.NewAuthService(pool)
	email := testutil.NewEmail()
	accountID := testutil.SeedAccount(t, pool, string(email), "Session User", "password")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	session, apiErr := svc.CreateSession(c, auth.AccountRow{ID: accountID})
	require.Nil(t, apiErr)
	require.NotNil(t, session.SessionToken)
	require.NotNil(t, session.ExpiresAt)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), *session.ExpiresAt, time.Minute)
}
