package handlers

import (
	"trenchcoat/internal/api"
	"trenchcoat/internal/api_error"
	"trenchcoat/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthServiceInterface interface {
	ValidateSignInCredentials(body api.SignInJSONRequestBody) []api.ErrorResponseDetail
	SignIn(c *gin.Context, body api.SignInJSONRequestBody) (*auth.SignInResponse, *api_error.ApiError)
	ValidateSignUpCredentials(body api.SignUpJSONRequestBody) []api.ErrorResponseDetail
	SignUp(c *gin.Context, body api.SignUpJSONRequestBody) (*auth.SignUpResponse, *api_error.ApiError)
	ParseAuthToken(authHeader string) (uuid.UUID, *api_error.ApiError)
	SignOut(c *gin.Context, tokenUUID uuid.UUID) *api_error.ApiError
}
