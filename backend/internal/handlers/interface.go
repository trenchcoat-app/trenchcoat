package handlers

import (
	"trenchcoat/internal/api"
	"trenchcoat/internal/dto/httperror"
	"trenchcoat/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthServiceInterface interface {
	ValidateSignInCredentials(body api.SignInJSONRequestBody) []api.ErrorResponseDetail
	SignIn(c *gin.Context, body api.SignInJSONRequestBody) (*auth.SignInResponse, *httperror.HttpError)
	ValidateSignUpCredentials(body api.SignUpJSONRequestBody) []api.ErrorResponseDetail
	SignUp(c *gin.Context, body api.SignUpJSONRequestBody) (*auth.SignUpResponse, *httperror.HttpError)
	ParseAuthToken(authHeader string) (uuid.UUID, *httperror.HttpError)
	SignOut(c *gin.Context, tokenUUID uuid.UUID) *httperror.HttpError
}
