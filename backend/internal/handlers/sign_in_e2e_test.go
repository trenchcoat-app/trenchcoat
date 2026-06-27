package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"trenchcoat/internal/api"
	"trenchcoat/internal/services/testutil"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
)

func TestE2E_SignIn_Success(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	router := testutil.SetupE2ERouter(t, pool)

	email := testutil.NewEmail()
	password := "secure-password"

	testutil.SeedAccount(t, pool, string(email), "E2E User", password)

	body := api.SignInJSONRequestBody{
		Email:    email,
		Password: password,
	}
	bodyBytes, _ := json.Marshal(body)
	w := testutil.PerformRequest(router, "POST", "/api/v1/auth/sign-in", bodyBytes)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp api.SignInOkResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, email, resp.Account.Email)
	assert.Equal(t, "E2E User", *resp.Account.DisplayName)
	assert.NotEqual(t, uuid.UUID{}, resp.Account.Id)

	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, "sid", cookies[0].Name)
	assert.NotEmpty(t, cookies[0].Value)
}
