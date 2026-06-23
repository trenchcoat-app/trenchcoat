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

	var errorDetails []struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	emailStr := strings.TrimSpace(string(body.Email))
	if emailStr == "" {
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

	if len(errorDetails) > 0 {
		c.JSON(http.StatusUnprocessableEntity, api.ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "Request validation failed.",
			Details: &errorDetails,
		})
		return
	}

	var (
		userID       openapi_types.UUID
		displayName  string
		passwordHash string
		status       string
	)

	err := s.DB.QueryRow(
		c.Request.Context(),
		"SELECT id, display_name, password_hash, status FROM account WHERE email = $1",
		strings.ToLower(emailStr),
	).Scan(&userID, &displayName, &passwordHash, &status)
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

	if status != "active" {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid email or password.",
			// Keep error message vague for security, but allow rejection if disabled
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, api.ErrorResponse{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid email or password.",
		})
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	var sessionToken openapi_types.UUID
	err = s.DB.QueryRow(
		c.Request.Context(),
		"INSERT INTO session (expires_at, ip_address, user_agent, account_id) VALUES ($1, $2, $3, $4) RETURNING token, expires_at",
		expiresAt,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		userID,
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
			Id:          userID,
			Email:       body.Email,
			DisplayName: &displayName,
		},
		Session: api.Session{
			Token:     sessionToken,
			ExpiresAt: expiresAt,
			AccountId: userID,
		},
	})
}
