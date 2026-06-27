package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"trenchcoat/internal/api"
	"trenchcoat/internal/dto/httperror"
	"trenchcoat/internal/handlers"
	"trenchcoat/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/mock"
)

var signInRoute = "/api/v1/auth/sign-in"

func TestSignIn_Success(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	token := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour)
	var email types.Email = "test@example.com"
	displayName := "Test User"

	mockAuth := NewMockAuthServiceInterface(t)
	mockAuth.EXPECT().
		ValidateSignInCredentials(mock.Anything).
		Return(nil).
		Once()
	mockAuth.EXPECT().
		SignIn(mock.Anything, mock.Anything).
		Return(&auth.SignInResponse{
			Account: &api.Account{
				Id:          uuid.New(),
				Email:       email,
				DisplayName: &displayName,
			},
			Session: &auth.Session{
				SessionToken: &token,
				ExpiresAt:    &expiresAt,
			},
		}, nil).
		Once()

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	bodyBytes, _ := json.Marshal(api.SignInJSONRequestBody{
		Email:    email,
		Password: "password123",
	})
	c.Request = httptest.NewRequest("POST", signInRoute, bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.SignIn(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp api.SignInOkResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, email, resp.Account.Email)
	assert.Equal(t, "Test User", *resp.Account.DisplayName)
	assert.NotEqual(t, uuid.Nil, resp.Account.Id)

	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, "sid", cookies[0].Name)
	assert.NotEmpty(t, cookies[0].Value)
	assert.Greater(t, cookies[0].MaxAge, 0)
}

func TestSignIn_BadJSON(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	mockAuth := NewMockAuthServiceInterface(t)

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", signInRoute, bytes.NewReader([]byte("{invalid")))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.SignIn(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "BAD_REQUEST", resp.Code)
	assert.Contains(t, resp.Message, "Malformed JSON payload")
}

func TestSignIn_ValidationError(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	mockAuth := NewMockAuthServiceInterface(t)
	mockAuth.EXPECT().
		ValidateSignInCredentials(mock.Anything).
		Return([]api.ErrorResponseDetail{
			{Field: "email", Message: "Email is required"},
			{Field: "password", Message: "Password is required"},
		})

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	bodyBytes, _ := json.Marshal(api.SignInJSONRequestBody{
		Email:    "test@example.com",
		Password: "",
	})
	c.Request = httptest.NewRequest("POST", signInRoute, bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.SignIn(c)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "VALIDATION_FAILED", resp.Code)
	require.NotNil(t, resp.Details)
	require.Len(t, *resp.Details, 2)
	assert.Equal(t, "email", (*resp.Details)[0].Field)
	assert.Equal(t, "password", (*resp.Details)[1].Field)
}

func TestSignIn_AuthError(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	mockAuth := NewMockAuthServiceInterface(t)
	mockAuth.EXPECT().
		ValidateSignInCredentials(mock.Anything).
		Return(nil)
	mockAuth.EXPECT().
		SignIn(mock.Anything, mock.Anything).
		Return(nil, httperror.SignInInvalidCredentialsError())

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	bodyBytes, _ := json.Marshal(api.SignInJSONRequestBody{
		Email:    "wrong@example.com",
		Password: "badpassword",
	})
	c.Request = httptest.NewRequest("POST", signInRoute, bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.SignIn(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", resp.Code)
	assert.Equal(t, "Invalid email or password.", resp.Message)
}
