package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"trenchcoat/internal/services/testUtil"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestE2E_SignOut_ViaCookie(t *testing.T) {
	t.Parallel()
	pool := testUtil.GetE2EPool(t)
	router := testUtil.SetupE2ERouter(t, pool)

	email := testUtil.NewEmail()
	accountID, err := testUtil.SeedAccount(pool, string(email), "E2E User", "secure-password")
	require.NoError(t, err)
	token, err := testUtil.SeedSession(pool, accountID)
	require.NoError(t, err)

	cookie := &http.Cookie{
		Name:  "sid",
		Value: token.String(),
	}
	w := testUtil.PerformRequest(router, "POST", "/api/v1/auth/sign-out", nil, cookie)

	assert.Equal(t, http.StatusNoContent, w.Code)

	respCookies := w.Result().Cookies()
	require.Len(t, respCookies, 1)
	assert.Equal(t, "sid", respCookies[0].Name)
	assert.Equal(t, "", respCookies[0].Value)
	assert.Less(t, respCookies[0].MaxAge, 0)
	assert.True(t, respCookies[0].Expires.Before(time.Now()))
}

func TestE2E_SignOut_ViaAuthHeader(t *testing.T) {
	t.Parallel()
	pool := testUtil.GetE2EPool(t)
	router := testUtil.SetupE2ERouter(t, pool)

	email := testUtil.NewEmail()
	accountID, err := testUtil.SeedAccount(pool, string(email), "E2E User", "secure-password")
	require.NoError(t, err)
	token, err := testUtil.SeedSession(pool, accountID)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/sign-out", nil)
	req.Header.Set("Authorization", "Bearer "+token.String())
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	respCookies := w.Result().Cookies()
	require.Len(t, respCookies, 1)
	assert.Equal(t, "sid", respCookies[0].Name)
	assert.Equal(t, "", respCookies[0].Value)
}
