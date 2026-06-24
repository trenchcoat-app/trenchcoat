package services

import (
	"strings"
	"time"
	"trenchcoat/internal/api"
	"trenchcoat/internal/api_error"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

type AccountRow struct {
	ID           openapi_types.UUID `db:"id"`
	DisplayName  string             `db:"display_name"`
	PasswordHash string             `db:"password_hash"`
	Status       string             `db:"status"`
}

func (auth *AuthService) SignIn(c *gin.Context, body api.SignInJSONRequestBody) (*api.Account, *api.Session, *api_error.ApiError) {
	account, apiErr := auth.getAccountRow(c, body)
	if apiErr != nil {
		return nil, nil, apiErr
	}

	sessionToken, expiresAt, apiErr := auth.createSession(c, *account)
	if apiErr != nil {
		return nil, nil, apiErr
	}

	return &api.Account{
			Id:          account.ID,
			Email:       body.Email,
			DisplayName: &account.DisplayName,
		},
		&api.Session{
			Token:     *sessionToken,
			ExpiresAt: *expiresAt,
			AccountId: account.ID,
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

func (auth *AuthService) getAccountRow(c *gin.Context, body api.SignInJSONRequestBody) (*AccountRow, *api_error.ApiError) {
	sql := `
		SELECT id, display_name, password_hash, status
		FROM account
		WHERE email = $1
	`

	rows, err := auth.DB.Query(
		c.Request.Context(),
		sql,
		strings.ToLower(string(body.Email)),
	)
	if err != nil {
		return nil, api_error.InternalServerError("Database query failed: " + err.Error())
	}

	account, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[AccountRow])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, api_error.SignInInvalidCredentialsError()
		}
		return nil, api_error.InternalServerError("Database query failed: " + err.Error())
	}

	if account.Status != "active" {
		return nil, api_error.SignInInvalidCredentialsError()
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(body.Password))
	if err != nil {
		return nil, api_error.SignInInvalidCredentialsError()
	}

	return &account, nil
}

func (auth *AuthService) createSession(c *gin.Context, account AccountRow) (sessionToken *openapi_types.UUID, expiresAt *time.Time, apiErr *api_error.ApiError) {
	expiresAt = auth.getNewSessionExpireTime()

	sql := `
		INSERT INTO session (expires_at, ip_address, user_agent, account_id)
		VALUES ($1, $2, $3, $4)
		RETURNING token, expires_at
	`

	err := auth.DB.QueryRow(
		c.Request.Context(),
		sql,
		expiresAt,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		account.ID,
	).Scan(&sessionToken, &expiresAt)

	if err != nil {
		apiErr = api_error.InternalServerError("Failed to create session: " + err.Error())
	}

	return
}

func (auth *AuthService) getNewSessionExpireTime() *time.Time {
	expireTime := time.Now().Add(24 * time.Hour)
	return &expireTime
}
