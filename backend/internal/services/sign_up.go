package services

import (
	"regexp"
	"strings"
	"trenchcoat/internal/api"
	"trenchcoat/internal/api_error"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(strings.ToLower(email))
}

func (auth *AuthService) ValidateSignUpCredentials(body api.SignUpJSONRequestBody) (errorDetails []api.ErrorResponseDetail) {
	nameTrimmed := strings.TrimSpace(body.Name)
	if nameTrimmed == "" {
		errorDetails = append(errorDetails, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{Field: "name", Message: "Name cannot be empty"})
	}

	emailStr := strings.TrimSpace(string(body.Email))
	if emailStr == "" || !isValidEmail(emailStr) {
		errorDetails = append(errorDetails, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{Field: "email", Message: "Invalid email format"})
	}

	if len(body.Password) < 8 {
		errorDetails = append(errorDetails, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{Field: "password", Message: "Password must be at least 8 characters long"})
	}

	return
}

func (auth *AuthService) SignUp(c *gin.Context, body api.SignUpJSONRequestBody) (*api.Account, *api.Session, *api_error.ApiError) {
	emailStr := strings.TrimSpace(string(body.Email))
	nameTrimmed := strings.TrimSpace(body.Name)

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
		return nil, nil, api_error.InternalServerError("Failed to query database: " + err.Error())
	}

	if exists {
		return nil, nil, api_error.SignUpEmailAlreadyExistsError()
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, api_error.InternalServerError("Failed to process password: " + err.Error())
	}

	userID, apiErr := auth.createAccount(c, emailStr, nameTrimmed, string(hashed))
	if apiErr != nil {
		return nil, nil, apiErr
	}

	var session *api.Session
	if body.AutoSignIn != nil && *body.AutoSignIn {
		sessionToken, expiresAt, apiErr := auth.createSession(c, AccountRow{ID: userID})
		if apiErr != nil {
			return nil, nil, apiErr
		}
		session = &api.Session{
			Token:     *sessionToken,
			ExpiresAt: *expiresAt,
			AccountId: userID,
		}
	}

	displayName := nameTrimmed
	return &api.Account{
			Id:          userID,
			Email:       openapi_types.Email(emailStr),
			DisplayName: &displayName,
		},
		session,
		nil
}

func (auth *AuthService) createAccount(c *gin.Context, email string, displayName string, passwordHash string) (openapi_types.UUID, *api_error.ApiError) {
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
