package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"trenchcoat/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) SignIn(c *gin.Context) {
	var body api.SignInJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Malformed JSON payload: " + err.Error(),
		})
		return
	}

	emailStr := strings.TrimSpace(string(body.Email))

	{
		validationErrorDetails := checkSignInValidationErrors(emailStr, body.Password)

		if len(validationErrorDetails) > 0 {
			c.JSON(http.StatusUnprocessableEntity, api.ErrorResponse{
				Code:    "VALIDATION_FAILED",
				Message: "Request validation failed.",
				Details: &validationErrorDetails,
			})
			return
		}
	}

	type AccountRow struct {
		ID           openapi_types.UUID `db:"id"`
		DisplayName  string             `db:"display_name"`
		PasswordHash string             `db:"password_hash"`
		Status       string             `db:"status"`
	}

	sql := `
            SELECT id, display_name, password_hash, status
            FROM account
            WHERE email = $1
        `

	rows, err := s.DB.Query(
		c.Request.Context(),
		sql,
		strings.ToLower(emailStr),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Database query failed: " + err.Error(),
		})
		return
	}

	account, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[AccountRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, api.ErrorResponse{
				Code:    "INVALID_CREDENTIALS",
				Message: "Invalid email or password.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Database query failed: " + err.Error(),
		})
		return
	}

	if account.Status != "active" {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid email or password.",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid email or password.",
		})
		return
	}

	// NOTE: Should probably be extracted elsewhere, this information feels too hidden inside the sign-in handler
	expiresAt := time.Now().Add(24 * time.Hour)

	var sessionToken openapi_types.UUID
	err = s.DB.QueryRow(
		c.Request.Context(),
		"INSERT INTO session (expires_at, ip_address, user_agent, account_id) VALUES ($1, $2, $3, $4) RETURNING token, expires_at",
		expiresAt,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		account.ID,
	).Scan(&sessionToken, &expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to create session: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, api.SignInOkResponse{
		Account: api.Account{
			Id:          account.ID,
			Email:       body.Email,
			DisplayName: &account.DisplayName,
		},
		Session: api.Session{
			Token:     sessionToken,
			ExpiresAt: expiresAt,
			AccountId: account.ID,
		},
	})
}

func checkSignInValidationErrors(email string, password string) (errorDetails []api.ErrorResponseDetail) {
	if email == "" {
		errorDetails = append(errorDetails, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{Field: "email", Message: "Email is required"})
	}

	if password == "" {
		errorDetails = append(errorDetails, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{Field: "password", Message: "Password is required"})
	}
	return
}
