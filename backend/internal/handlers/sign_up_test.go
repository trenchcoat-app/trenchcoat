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

var signUpRoute = "/api/v1/auth/sign-up"

func TestSignUp_SuccessWithoutAutoSignIn(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	var email types.Email = "new@example.com"
	displayName := "New User"
	password := "password123"

	mockAuth := NewMockAuthServiceInterface(t)
	mockAuth.EXPECT().
		ValidateSignUpCredentials(mock.Anything).
		Return(nil)
	mockAuth.EXPECT().
		SignUp(mock.Anything, mock.Anything).
		Return(&auth.SignUpResponse{
			Account: &api.Account{
				Id:          userID,
				Email:       email,
				DisplayName: &displayName,
			},
			Session: nil,
		}, nil)

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	bodyBytes, _ := json.Marshal(api.SignUpJSONRequestBody{
		Email:       email,
		Password:    password,
		DisplayName: displayName,
	})
	c.Request = httptest.NewRequest("POST", signUpRoute, bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.SignUp(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp api.SignUpOkResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, email, resp.Account.Email)
	assert.Equal(t, displayName, *resp.Account.DisplayName)
	assert.Equal(t, userID, resp.Account.Id)

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 0, "no cookie should be set when AutoSignIn is false")
}

func TestSignUp_SuccessWithAutoSignIn(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	token := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour)
	var email types.Email = "new@example.com"
	displayName := "New User"
	password := "password123"

	mockAuth := NewMockAuthServiceInterface(t)
	mockAuth.EXPECT().
		ValidateSignUpCredentials(mock.Anything).
		Return(nil)
	mockAuth.EXPECT().
		SignUp(mock.Anything, mock.Anything).
		Return(&auth.SignUpResponse{
			Account: &api.Account{
				Id:          uuid.New(),
				Email:       email,
				DisplayName: &displayName,
			},
			Session: &auth.Session{
				SessionToken: &token,
				ExpiresAt:    &expiresAt,
			},
		}, nil)

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	autoSignIn := true
	bodyBytes, _ := json.Marshal(api.SignUpJSONRequestBody{
		Email:       email,
		Password:    password,
		DisplayName: displayName,
		AutoSignIn:  &autoSignIn,
	})
	c.Request = httptest.NewRequest("POST", signUpRoute, bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.SignUp(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	cookies := w.Result().Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, "sid", cookies[0].Name)
	assert.Equal(t, token.String(), cookies[0].Value)
	assert.Greater(t, cookies[0].MaxAge, 0)
}

func TestSignUp_BadJSON(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	mockAuth := NewMockAuthServiceInterface(t)

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("POST", signUpRoute, bytes.NewReader([]byte("{invalid")))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.SignUp(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "BAD_REQUEST", resp.Code)
	assert.Contains(t, resp.Message, "Malformed JSON payload")
}

func TestSignUp_ValidationError(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	mockAuth := NewMockAuthServiceInterface(t)
	mockAuth.EXPECT().
		ValidateSignUpCredentials(mock.Anything).
		Return([]api.ErrorResponseDetail{
			{Field: "name", Message: "Name cannot be empty"},
			{Field: "password", Message: "Password must be at least 8 characters long"},
		})

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	bodyBytes, _ := json.Marshal(api.SignUpJSONRequestBody{
		Email:       "new@example.com",
		Password:    "short",
		DisplayName: "",
	})
	c.Request = httptest.NewRequest("POST", signUpRoute, bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.SignUp(c)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "VALIDATION_FAILED", resp.Code)
	require.NotNil(t, resp.Details)
	require.Len(t, *resp.Details, 2)
}

func TestSignUp_EmailAlreadyExists(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	mockAuth := NewMockAuthServiceInterface(t)
	mockAuth.EXPECT().
		ValidateSignUpCredentials(mock.Anything).
		Return(nil)
	mockAuth.EXPECT().
		SignUp(mock.Anything, mock.Anything).
		Return(nil, httperror.SignUpEmailAlreadyExistsError())

	srv := handlers.NewServer(mockAuth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	bodyBytes, _ := json.Marshal(api.SignUpJSONRequestBody{
		Email:       "existing@example.com",
		Password:    "password123",
		DisplayName: "Existing User",
	})
	c.Request = httptest.NewRequest("POST", signUpRoute, bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.SignUp(c)

	assert.Equal(t, http.StatusConflict, w.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "EMAIL_ALREADY_EXISTS", resp.Code)
	assert.Equal(t, "An account with this email address already exists.", resp.Message)
}
