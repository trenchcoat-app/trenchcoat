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

func TestE2E_SignUp_Success(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	router := testutil.SetupE2ERouter(t, pool)

	autoSignIn := false
	body := api.SignUpJSONRequestBody{
		Email:       testutil.NewEmail(),
		Password:    "secure-password",
		DisplayName: "E2E User",
		AutoSignIn:  &autoSignIn,
	}
	bodyBytes, _ := json.Marshal(body)
	w := testutil.PerformRequest(router, "POST", "/api/v1/auth/sign-up", bodyBytes)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp api.SignUpOkResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, body.Email, resp.Account.Email)
	assert.Equal(t, "E2E User", *resp.Account.DisplayName)
	assert.NotEqual(t, uuid.UUID{}, resp.Account.Id)

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 0)
}

func TestE2E_SignUp_SuccessWithAutoSignIn(t *testing.T) {
	t.Parallel()
	pool := testutil.GetE2EPool(t)
	router := testutil.SetupE2ERouter(t, pool)

	autoSignIn := true
	body := api.SignUpJSONRequestBody{
		Email:       testutil.NewEmail(),
		Password:    "secure-password",
		DisplayName: "E2E User",
		AutoSignIn:  &autoSignIn,
	}
	bodyBytes, _ := json.Marshal(body)
	w := testutil.PerformRequest(router, "POST", "/api/v1/auth/sign-up", bodyBytes)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp api.SignUpOkResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, body.Email, resp.Account.Email)
	assert.Equal(t, "E2E User", *resp.Account.DisplayName)

	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, "sid", cookies[0].Name)
	assert.NotEmpty(t, cookies[0].Value)
}
