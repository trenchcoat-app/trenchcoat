package auth_test

import (
	"testing"

	"trenchcoat/internal/services/auth"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
)

func TestParseAuthToken_Empty(t *testing.T) {
	svc := auth.NewAuthService(nil)

	_, httpErr := svc.ParseAuthToken("")
	require.NotNil(t, httpErr)
	assert.Equal(t, "UNAUTHORIZED", httpErr.Code)
	assert.Equal(t, "Missing authorization token.", httpErr.Message)
}

func TestParseAuthToken_BearerPrefix(t *testing.T) {
	svc := auth.NewAuthService(nil)

	token := uuid.New()
	parsed, httpErr := svc.ParseAuthToken("Bearer " + token.String())
	require.Nil(t, httpErr)
	assert.Equal(t, token, parsed)
}

func TestParseAuthToken_BearerLowercase(t *testing.T) {
	svc := auth.NewAuthService(nil)

	token := uuid.New()
	parsed, httpErr := svc.ParseAuthToken("bearer " + token.String())
	require.Nil(t, httpErr)
	assert.Equal(t, token, parsed)
}

func TestParseAuthToken_RawUUID(t *testing.T) {
	svc := auth.NewAuthService(nil)

	token := uuid.New()
	parsed, httpErr := svc.ParseAuthToken(token.String())
	require.Nil(t, httpErr)
	assert.Equal(t, token, parsed)
}

func TestParseAuthToken_InvalidUUID(t *testing.T) {
	svc := auth.NewAuthService(nil)

	_, httpErr := svc.ParseAuthToken("not-a-uuid")
	require.NotNil(t, httpErr)
	assert.Equal(t, "UNAUTHORIZED", httpErr.Code)
	assert.Equal(t, "Invalid authorization token format.", httpErr.Message)
}

func TestParseAuthToken_WhitespaceAfterBearer(t *testing.T) {
	svc := auth.NewAuthService(nil)

	token := uuid.New()
	parsed, httpErr := svc.ParseAuthToken("Bearer   " + token.String())
	require.Nil(t, httpErr)
	assert.Equal(t, token, parsed)
}

func TestParseAuthToken_BearerOnly(t *testing.T) {
	svc := auth.NewAuthService(nil)

	_, httpErr := svc.ParseAuthToken("Bearer ")
	require.NotNil(t, httpErr)
	assert.Equal(t, "UNAUTHORIZED", httpErr.Code)
	assert.Equal(t, "Missing authorization token.", httpErr.Message)
}
