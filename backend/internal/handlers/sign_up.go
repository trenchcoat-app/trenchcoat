package handlers

import (
	"net/http"
	"regexp"
	"strings"
	"time"
	"trenchcoat/internal/api"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(strings.ToLower(email))
}

func (s *Server) SignUp(c *gin.Context) {
	var body api.SignUpJSONRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Code:    "BAD_REQUEST",
			Message: "Malformed JSON payload: " + err.Error(),
		})
		return
	}

	var details []struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	nameTrimmed := strings.TrimSpace(body.Name)
	if nameTrimmed == "" {
		details = append(details, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{Field: "name", Message: "Name cannot be empty"})
	}

	emailStr := strings.TrimSpace(string(body.Email))
	if emailStr == "" || !isValidEmail(emailStr) {
		details = append(details, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{Field: "email", Message: "Invalid email format"})
	}

	if len(body.Password) < 8 {
		details = append(details, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{Field: "password", Message: "Password must be at least 8 characters long"})
	}

	if len(details) > 0 {
		c.JSON(http.StatusUnprocessableEntity, api.ErrorResponse{
			Code:    "VALIDATION_FAILED",
			Message: "Request validation failed.",
			Details: &details,
		})
		return
	}

	var exists bool
	err := s.DB.QueryRow(
		c.Request.Context(),
		"SELECT EXISTS(SELECT 1 FROM account WHERE email = $1)",
		strings.ToLower(emailStr),
	).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to query database: " + err.Error(),
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, api.ErrorResponse{
			Code:    "EMAIL_ALREADY_EXISTS",
			Message: "An account with this email address already exists.",
		})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to process password: " + err.Error(),
		})
		return
	}

	var userID openapi_types.UUID
	err = s.DB.QueryRow(
		c.Request.Context(),
		"INSERT INTO account (email, display_name, password_hash, status) VALUES ($1, $2, $3, 'active') RETURNING id",
		strings.ToLower(emailStr),
		nameTrimmed,
		string(hashed),
	).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to create account: " + err.Error(),
		})
		return
	}

	var session *api.Session
	if body.AutoSignIn != nil && *body.AutoSignIn {
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
				Message: "Account created but failed to start session: " + err.Error(),
			})
			return
		}
		session = &api.Session{
			Token:     sessionToken,
			ExpiresAt: expiresAt,
			AccountId: userID,
		}
	}

	displayName := nameTrimmed
	c.JSON(http.StatusCreated, api.SignUpOkResponse{
		Account: api.Account{
			Id:          userID,
			Email:       openapi_types.Email(emailStr),
			DisplayName: &displayName,
		},
		Session: session,
	})
}
