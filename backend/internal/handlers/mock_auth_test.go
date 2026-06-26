package handlers_test

import (
	"trenchcoat/internal/api"
	"trenchcoat/internal/api_error"
	"trenchcoat/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockAuthService struct {
	validateSignInCredentialsFn func(body api.SignInJSONRequestBody) []api.ErrorResponseDetail
	signInFn                    func(c *gin.Context, body api.SignInJSONRequestBody) (*auth.SignInResponse, *api_error.ApiError)
	validateSignUpCredentialsFn func(body api.SignUpJSONRequestBody) []api.ErrorResponseDetail
	signUpFn                    func(c *gin.Context, body api.SignUpJSONRequestBody) (*auth.SignUpResponse, *api_error.ApiError)
	parseAuthTokenFn            func(authHeader string) (uuid.UUID, *api_error.ApiError)
	signOutFn                   func(c *gin.Context, tokenUUID uuid.UUID) *api_error.ApiError
}

func (m *mockAuthService) ValidateSignInCredentials(body api.SignInJSONRequestBody) []api.ErrorResponseDetail {
	return m.validateSignInCredentialsFn(body)
}

func (m *mockAuthService) SignIn(c *gin.Context, body api.SignInJSONRequestBody) (*auth.SignInResponse, *api_error.ApiError) {
	return m.signInFn(c, body)
}

func (m *mockAuthService) ValidateSignUpCredentials(body api.SignUpJSONRequestBody) []api.ErrorResponseDetail {
	return m.validateSignUpCredentialsFn(body)
}

func (m *mockAuthService) SignUp(c *gin.Context, body api.SignUpJSONRequestBody) (*auth.SignUpResponse, *api_error.ApiError) {
	return m.signUpFn(c, body)
}

func (m *mockAuthService) ParseAuthToken(authHeader string) (uuid.UUID, *api_error.ApiError) {
	return m.parseAuthTokenFn(authHeader)
}

func (m *mockAuthService) SignOut(c *gin.Context, tokenUUID uuid.UUID) *api_error.ApiError {
	return m.signOutFn(c, tokenUUID)
}
