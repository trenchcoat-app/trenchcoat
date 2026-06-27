package testutil

import (
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}
