package auth_test

import (
	"testing"

	"trenchcoat/internal/api"
	"trenchcoat/internal/services/auth"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/oapi-codegen/runtime/types"
)

func TestValidateSignUpCredentials(t *testing.T) {
	svc := auth.NewAuthService(nil)

	tests := []struct {
		name        string
		email       types.Email
		password    string
		displayName string
		expected    []api.ErrorResponseDetail
	}{
		{
			name:        "empty name",
			email:       "test@example.com",
			password:    "password123",
			displayName: "",
			expected:    []api.ErrorResponseDetail{{Field: "name", Message: "Name cannot be empty"}},
		},
		{
			name:        "empty name after trim",
			email:       "test@example.com",
			password:    "password123",
			displayName: "   ",
			expected:    []api.ErrorResponseDetail{{Field: "name"}},
		},
		{
			name:        "short password",
			email:       "test@example.com",
			password:    "1234567",
			displayName: "Test User",
			expected:    []api.ErrorResponseDetail{{Field: "password", Message: "Password must be at least 8 characters long"}},
		},
		{
			name:        "valid",
			email:       "test@example.com",
			password:    "password123",
			displayName: "Test User",
			expected:    nil,
		},
		{
			name:        "both invalid",
			email:       "test@example.com",
			password:    "short",
			displayName: "",
			expected:    []api.ErrorResponseDetail{{Field: "name"}, {Field: "password"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			details := svc.ValidateSignUpCredentials(api.SignUpJSONRequestBody{
				Email:       tt.email,
				Password:    tt.password,
				DisplayName: tt.displayName,
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
