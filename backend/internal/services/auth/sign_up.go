package auth

import (
	"strings"
	"trenchcoat/internal/api"
	"trenchcoat/internal/api_error"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

type SignUpResponse struct {
	Account *api.Account
	Session *Session
}

func (auth *AuthService) ValidateSignUpCredentials(body api.SignUpJSONRequestBody) (errorDetails []api.ErrorResponseDetail) {
	nameTrimmed := strings.TrimSpace(body.DisplayName)
	if nameTrimmed == "" {
		errorDetails = append(errorDetails, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{Field: "name", Message: "Name cannot be empty"})
	}

	if len(body.Password) < 8 {
		errorDetails = append(errorDetails, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{Field: "password", Message: "Password must be at least 8 characters long"})
	}

	return
}

func (auth *AuthService) SignUp(c *gin.Context, body api.SignUpJSONRequestBody) (*SignUpResponse, *api_error.ApiError) {
	emailStr := strings.TrimSpace(string(body.Email))
	nameTrimmed := strings.TrimSpace(body.DisplayName)

	sqlExists := `
		SELECT EXISTS(SELECT 1 FROM account WHERE email = $1)
	`

	var exists bool
	err := auth.DB.QueryRow(
		c.Request.Context(),
		sqlExists,
		strings.ToLower(emailStr),
	).Scan(&exists)
	if err != nil {
		return nil, api_error.InternalServerError("Failed to query database: " + err.Error())
	}

	if exists {
		return nil, api_error.SignUpEmailAlreadyExistsError()
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, api_error.InternalServerError("Failed to process password: " + err.Error())
	}

	userID, apiErr := auth.CreateAccount(c, emailStr, nameTrimmed, string(hashed))
	if apiErr != nil {
		return nil, apiErr
	}

	var session Session
	if body.AutoSignIn != nil && *body.AutoSignIn {
		var apiErr *api_error.ApiError
		session, apiErr = auth.CreateSession(c, accountRow{ID: userID})
		if apiErr != nil {
			return nil, apiErr
		}
	}

	displayName := nameTrimmed
	return &SignUpResponse{
		&api.Account{
			Id:          userID,
			Email:       openapi_types.Email(emailStr),
			DisplayName: &displayName,
		},
		&session,
	}, nil
}

func (auth *AuthService) CreateAccount(c *gin.Context, email string, displayName string, passwordHash string) (openapi_types.UUID, *api_error.ApiError) {
	var userID openapi_types.UUID

	sql := `
		INSERT INTO account (email, display_name, password_hash, status)
		VALUES ($1, $2, $3, 'active')
		RETURNING id
	`

	err := auth.DB.QueryRow(
		c.Request.Context(),
		sql,
		strings.ToLower(email),
		displayName,
		passwordHash,
	).Scan(&userID)
	if err != nil {
		return userID, api_error.InternalServerError("Failed to create account: " + err.Error())
	}
	return userID, nil
}
