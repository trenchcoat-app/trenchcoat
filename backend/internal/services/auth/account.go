package auth

import (
	"strings"
	"trenchcoat/internal/api"
	"trenchcoat/internal/httperror"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func (auth *AuthService) GetAccountRow(c *gin.Context, body api.SignInJSONRequestBody) (*AccountRow, *httperror.HttpError) {
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
		return nil, httperror.InternalServerError("Database query failed: " + err.Error())
	}

	account, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[AccountRow])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, httperror.SignInInvalidCredentialsError()
		}
		return nil, httperror.InternalServerError("Database query failed: " + err.Error())
	}

	if account.Status != "active" {
		return nil, httperror.SignInInvalidCredentialsError()
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(body.Password))
	if err != nil {
		return nil, httperror.SignInInvalidCredentialsError()
	}

	return &account, nil
}
