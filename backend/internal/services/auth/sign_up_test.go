package auth_test

import (
	"testing"

	"trenchcoat/internal/api"
	"trenchcoat/internal/services/auth"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestValidateSignUpCredentials_EmptyName(t *testing.T) {
	svc := auth.NewAuthService(nil)

	details := svc.ValidateSignUpCredentials(api.SignUpJSONRequestBody{
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "",
	})

	require.Len(t, details, 1)
	assert.Equal(t, "name", details[0].Field)
	assert.Equal(t, "Name cannot be empty", details[0].Message)
}

func TestValidateSignUpCredentials_EmptyNameAfterTrim(t *testing.T) {
	svc := auth.NewAuthService(nil)

	details := svc.ValidateSignUpCredentials(api.SignUpJSONRequestBody{
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "   ",
	})

	require.Len(t, details, 1)
	assert.Equal(t, "name", details[0].Field)
}

func TestValidateSignUpCredentials_ShortPassword(t *testing.T) {
	svc := auth.NewAuthService(nil)

	details := svc.ValidateSignUpCredentials(api.SignUpJSONRequestBody{
		Email:       "test@example.com",
		Password:    "1234567",
		DisplayName: "Test User",
	})

	require.Len(t, details, 1)
	assert.Equal(t, "password", details[0].Field)
	assert.Equal(t, "Password must be at least 8 characters long", details[0].Message)
}

func TestValidateSignUpCredentials_Valid(t *testing.T) {
	svc := auth.NewAuthService(nil)

	details := svc.ValidateSignUpCredentials(api.SignUpJSONRequestBody{
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	})

	assert.Len(t, details, 0)
}

func TestValidateSignUpCredentials_BothInvalid(t *testing.T) {
	svc := auth.NewAuthService(nil)

	details := svc.ValidateSignUpCredentials(api.SignUpJSONRequestBody{
		Email:       "test@example.com",
		Password:    "short",
		DisplayName: "",
	})

	require.Len(t, details, 2)
	assert.Equal(t, "name", details[0].Field)
	assert.Equal(t, "password", details[1].Field)
}
