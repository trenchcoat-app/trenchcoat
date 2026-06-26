package auth_test

import (
	"testing"

	"trenchcoat/internal/api"
	"trenchcoat/internal/services/auth"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestValidateSignInCredentials_EmptyEmail(t *testing.T) {
	svc := auth.NewAuthService(nil)

	details := svc.ValidateSignInCredentials(api.SignInJSONRequestBody{
		Email:    "",
		Password: "password123",
	})

	require.Len(t, details, 1)
	assert.Equal(t, "email", details[0].Field)
	assert.Equal(t, "Email is required", details[0].Message)
}

func TestValidateSignInCredentials_EmptyPassword(t *testing.T) {
	svc := auth.NewAuthService(nil)

	details := svc.ValidateSignInCredentials(api.SignInJSONRequestBody{
		Email:    "test@example.com",
		Password: "",
	})

	require.Len(t, details, 1)
	assert.Equal(t, "password", details[0].Field)
	assert.Equal(t, "Password is required", details[0].Message)
}

func TestValidateSignInCredentials_BothEmpty(t *testing.T) {
	svc := auth.NewAuthService(nil)

	details := svc.ValidateSignInCredentials(api.SignInJSONRequestBody{
		Email:    "",
		Password: "",
	})

	require.Len(t, details, 2)
	assert.Equal(t, "email", details[0].Field)
	assert.Equal(t, "password", details[1].Field)
}

func TestValidateSignInCredentials_Valid(t *testing.T) {
	svc := auth.NewAuthService(nil)

	details := svc.ValidateSignInCredentials(api.SignInJSONRequestBody{
		Email:    "test@example.com",
		Password: "password123",
	})

	assert.Len(t, details, 0)
}
