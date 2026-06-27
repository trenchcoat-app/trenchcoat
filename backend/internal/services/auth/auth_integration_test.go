package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"

	"trenchcoat/internal/api"
	"trenchcoat/internal/db/dbtest"
	"trenchcoat/internal/services/auth"
)

var (
	testPool     *pgxpool.Pool
	testPoolOnce sync.Once
)

func getTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	testPoolOnce.Do(func() {
		var err error
		testPool, _, err = dbtest.SetupDBMain()
		require.NoError(t, err, "Failed to set up test database")
	})
	return testPool
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}

func seedAccount(t *testing.T, pool *pgxpool.Pool, email, displayName, password string) types.UUID {
	t.Helper()
	hash := hashPassword(t, password)
	var id types.UUID
	err := pool.QueryRow(
		context.Background(),
		`INSERT INTO account (email, display_name, password_hash, status) VALUES ($1, $2, $3, 'active') RETURNING id`,
		email, displayName, hash,
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func seedSession(t *testing.T, pool *pgxpool.Pool, accountID types.UUID) (types.UUID, time.Time) {
	t.Helper()
	var token types.UUID
	var expiresAt time.Time
	err := pool.QueryRow(
		context.Background(),
		`INSERT INTO session (expires_at, account_id) VALUES ($1, $2) RETURNING token, expires_at`,
		time.Now().Add(24*time.Hour), accountID,
	).Scan(&token, &expiresAt)
	require.NoError(t, err)
	return token, expiresAt
}

func newEmail() types.Email {
	return types.Email("test-" + uuid.NewString() + "@example.com")
}

// GetAccountRow tests

func TestGetAccountRow_Success(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)
	email := newEmail()
	seedAccount(t, pool, string(email), "Test User", "correct-password")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := api.SignInJSONRequestBody{
		Email:    email,
		Password: "correct-password",
	}

	account, apiErr := svc.GetAccountRow(c, body)
	require.Nil(t, apiErr)
	require.NotNil(t, account)
	assert.Equal(t, "Test User", account.DisplayName)
	assert.Equal(t, "active", account.Status)
}

func TestGetAccountRow_InvalidCredentials(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)

	t.Run("wrong password", func(t *testing.T) {
		t.Parallel()
		email := newEmail()
		seedAccount(t, pool, string(email), "Test User", "correct-password")

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

		body := api.SignInJSONRequestBody{
			Email:    email,
			Password: "wrong-password",
		}
		_, apiErr := svc.GetAccountRow(c, body)
		require.NotNil(t, apiErr)
		assert.Equal(t, http.StatusUnauthorized, apiErr.Status)
	})

	t.Run("unknown email", func(t *testing.T) {
		t.Parallel()
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

		body := api.SignInJSONRequestBody{
			Email:    newEmail(),
			Password: "any-password",
		}
		_, apiErr := svc.GetAccountRow(c, body)
		require.NotNil(t, apiErr)
		assert.Equal(t, http.StatusUnauthorized, apiErr.Status)
	})
}

func TestGetAccountRow_DisabledAccount(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)
	email := newEmail()
	hash := hashPassword(t, "any-password")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	var id types.UUID
	err := pool.QueryRow(
		context.Background(),
		`INSERT INTO account (email, display_name, password_hash, status) VALUES ($1, $2, $3, 'disabled') RETURNING id`,
		string(email), "Disabled User", hash,
	).Scan(&id)
	require.NoError(t, err)

	body := api.SignInJSONRequestBody{
		Email:    email,
		Password: "any-password",
	}

	_, apiErr := svc.GetAccountRow(c, body)
	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.Status)
}

// CreateAccount tests

func TestCreateAccount_Success(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)
	email := newEmail()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	hash := hashPassword(t, "secure-password")
	id, apiErr := svc.CreateAccount(c, string(email), "New User", hash)
	require.Nil(t, apiErr)

	var storedEmail string
	var displayName string
	err := pool.QueryRow(
		context.Background(),
		`SELECT email, display_name FROM account WHERE id = $1`, id,
	).Scan(&storedEmail, &displayName)
	require.NoError(t, err)
	assert.Equal(t, string(email), storedEmail)
	assert.Equal(t, "New User", displayName)
}

func TestCreateAccount_DuplicateEmail(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)
	email := newEmail()
	seedAccount(t, pool, string(email), "First User", "password1")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	hash := hashPassword(t, "password2")
	_, apiErr := svc.CreateAccount(c, string(email), "Second User", hash)
	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
}

// CreateSession tests

func TestCreateSession_Success(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)
	email := newEmail()
	accountID := seedAccount(t, pool, string(email), "Session User", "password")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	session, apiErr := svc.CreateSession(c, auth.AccountRow{ID: accountID})
	require.Nil(t, apiErr)
	require.NotNil(t, session.SessionToken)
	require.NotNil(t, session.ExpiresAt)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), *session.ExpiresAt, time.Minute)
}

// SignIn tests

func TestSignIn_Success(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)
	email := newEmail()
	password := "the-password"
	seedAccount(t, pool, string(email), "Signin User", password)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := api.SignInJSONRequestBody{
		Email:    email,
		Password: password,
	}

	resp, apiErr := svc.SignIn(c, body)
	require.Nil(t, apiErr)
	require.NotNil(t, resp)
	assert.Equal(t, "Signin User", *resp.Account.DisplayName)
	assert.Equal(t, email, resp.Account.Email)
	require.NotNil(t, resp.Session.SessionToken)
}

func TestSignIn_InvalidCredentials(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)
	email := newEmail()
	seedAccount(t, pool, string(email), "Test User", "correct-password")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := api.SignInJSONRequestBody{
		Email:    email,
		Password: "wrong-password",
	}

	_, apiErr := svc.SignIn(c, body)
	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.Status)
}

func TestSignIn_DisabledAccount(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)
	email := newEmail()
	hash := hashPassword(t, "password")

	var id types.UUID
	err := pool.QueryRow(
		context.Background(),
		`INSERT INTO account (email, display_name, password_hash, status) VALUES ($1, $2, $3, 'disabled') RETURNING id`,
		string(email), "Disabled User", hash,
	).Scan(&id)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := api.SignInJSONRequestBody{
		Email:    email,
		Password: "password",
	}

	_, apiErr := svc.SignIn(c, body)
	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.Status)
}

// SignUp tests

func TestSignUp_Success(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)
	email := newEmail()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	autoSignIn := false
	body := api.SignUpJSONRequestBody{
		Email:       email,
		Password:    "secure-password",
		DisplayName: "New User",
		AutoSignIn:  &autoSignIn,
	}

	resp, apiErr := svc.SignUp(c, body)
	require.Nil(t, apiErr)
	require.NotNil(t, resp)
	assert.Equal(t, "New User", *resp.Account.DisplayName)
	assert.Equal(t, email, resp.Account.Email)
}

func TestSignUp_DuplicateEmail(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)
	email := newEmail()
	seedAccount(t, pool, string(email), "First User", "password1")

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	body := api.SignUpJSONRequestBody{
		Email:       email,
		Password:    "another-password",
		DisplayName: "Second User",
	}

	_, apiErr := svc.SignUp(c, body)
	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusConflict, apiErr.Status)
	assert.Equal(t, "EMAIL_ALREADY_EXISTS", apiErr.Code)
}

func TestSignUp_AutoSignIn(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)
	email := newEmail()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	autoSignIn := true
	body := api.SignUpJSONRequestBody{
		Email:       email,
		Password:    "secure-password",
		DisplayName: "Auto User",
		AutoSignIn:  &autoSignIn,
	}

	resp, apiErr := svc.SignUp(c, body)
	require.Nil(t, apiErr)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Session)
	require.NotNil(t, resp.Session.SessionToken)
}

// SignOut tests

func TestSignOut_Success(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)
	email := newEmail()
	accountID := seedAccount(t, pool, string(email), "Signout User", "password")
	token, _ := seedSession(t, pool, accountID)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	apiErr := svc.SignOut(c, token)
	require.Nil(t, apiErr)

	var count int
	err := pool.QueryRow(
		context.Background(), `SELECT COUNT(*) FROM session WHERE token = $1`, token,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestSignOut_NotFound(t *testing.T) {
	t.Parallel()
	pool := getTestPool(t)
	svc := auth.NewAuthService(pool)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	token := uuid.New()
	apiErr := svc.SignOut(c, token)
	require.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.Status)
	assert.Equal(t, "UNAUTHORIZED", apiErr.Code)
}
