package auth_test

import (
	"testing"

	"trenchcoat/internal/api"
	"trenchcoat/internal/services/auth"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/oapi-codegen/runtime/types"
)

func TestValidateSignInCredentials(t *testing.T) {
	svc := auth.NewAuthService(nil)

	tests := []struct {
		name     string
		email    types.Email
		password string
		expected []api.ErrorResponseDetail
	}{
		{
			name:     "empty email",
			email:    "",
			password: "password123",
			expected: []api.ErrorResponseDetail{{Field: "email", Message: "Email is required"}},
		},
		{
			name:     "empty password",
			email:    "test@example.com",
			password: "",
			expected: []api.ErrorResponseDetail{{Field: "password", Message: "Password is required"}},
		},
		{
			name:     "both empty",
			email:    "",
			password: "",
			expected: []api.ErrorResponseDetail{{Field: "email"}, {Field: "password"}},
		},
		{
			name:     "valid",
			email:    "test@example.com",
			password: "password123",
			},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			details := svc.ValidateSignInCredentials(api.SignInJSONRequestBody{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Len(t, details, len(tt.expected))
			for i, e := range tt.expected {
				assert.Equal(t, e.Field, details[i].Field)
				if e.Message != "" {
					assert.Equal(t, e.Message, details[i].Message)
				}
			}
		})
	}
}
