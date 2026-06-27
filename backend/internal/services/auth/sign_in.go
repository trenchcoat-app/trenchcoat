package auth

import (
	"trenchcoat/internal/api"
	"trenchcoat/internal/api_error"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type accountRow struct {
	ID           openapi_types.UUID `db:"id"`
	DisplayName  string             `db:"display_name"`
	PasswordHash string             `db:"password_hash"`
	Status       string             `db:"status"`
}

type SignInResponse struct {
	Account *api.Account
	Session *Session
}

func (auth *AuthService) SignIn(c *gin.Context, body api.SignInJSONRequestBody) (*SignInResponse, *api_error.ApiError) {
	account, apiErr := auth.GetAccountRow(c, body)
	if apiErr != nil {
		return nil, apiErr
	}

	session, apiErr := auth.CreateSession(c, *account)
	if apiErr != nil {
		return nil, apiErr
	}

	return &SignInResponse{
			&api.Account{
				Id:          account.ID,
				Email:       body.Email,
				DisplayName: &account.DisplayName,
			},
			&Session{
				SessionToken: session.SessionToken,
				ExpiresAt:    session.ExpiresAt,
			},
		},
		nil
}

func (auth *AuthService) ValidateSignInCredentials(body api.SignInJSONRequestBody) (errorDetails []api.ErrorResponseDetail) {
	if body.Email == "" {
		errorDetails = append(errorDetails, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{Field: "email", Message: "Email is required"})
	}

	if body.Password == "" {
		errorDetails = append(errorDetails, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{Field: "password", Message: "Password is required"})
	}
	return
}
