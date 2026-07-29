package auth

import (
	"strings"
	"trenchcoat/internal/api"
	"trenchcoat/internal/dto/httperror"

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
		errorDetails = append(errorDetails, api.ErrorResponseDetail{Field: "name", Message: "Name cannot be empty"})
	}

	if len(body.Password) < 8 {
		errorDetails = append(errorDetails, api.ErrorResponseDetail{Field: "password", Message: "Password must be at least 8 characters long"})
	}

	return
}

func (auth *AuthService) SignUp(c *gin.Context, body api.SignUpJSONRequestBody) (*SignUpResponse, *httperror.HTTPError) {
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
		return nil, httperror.InternalServerError("Failed to query database: " + err.Error())
	}

	if exists {
		return nil, httperror.SignUpEmailAlreadyExistsError()
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, httperror.InternalServerError("Failed to process password: " + err.Error())
	}

	userID, httpErr := auth.CreateAccount(c, emailStr, nameTrimmed, string(hashed))
	if httpErr != nil {
		return nil, httpErr
	}

	var session Session
	if body.AutoSignIn != nil && *body.AutoSignIn {
		session, httpErr = auth.CreateSession(c, AccountRow{ID: userID})
		if httpErr != nil {
			return nil, httpErr
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

func (auth *AuthService) CreateAccount(c *gin.Context, email string, displayName string, passwordHash string) (openapi_types.UUID, *httperror.HTTPError) {
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
		return userID, httperror.InternalServerError("Failed to create account: " + err.Error())
	}
	return userID, nil
}
