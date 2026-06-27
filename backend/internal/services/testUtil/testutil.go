package testUtil

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"trenchcoat/internal/api"
	"trenchcoat/internal/db/dbtest"
	"trenchcoat/internal/handlers"
	"trenchcoat/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

var (
	e2ePool     *pgxpool.Pool
	e2ePoolOnce sync.Once
)

func GetE2EPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	e2ePoolOnce.Do(func() {
		var err error
		e2ePool, _, err = dbtest.SetupDBMain()
		require.NoError(t, err, "Failed to set up test database")
	})
	return e2ePool
}

func SetupE2ERouter(t *testing.T, pool *pgxpool.Pool) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	authService := auth.NewAuthService(pool)
	srv := handlers.NewServer(authService)
	router := gin.New()
	api.RegisterHandlers(router, srv)
	return router
}

func PerformRequest(router *gin.Engine, method, path string, body []byte, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	router.ServeHTTP(w, req)
	return w
}

func NewEmail() openapi_types.Email {
	return openapi_types.Email("e2e-test-" + uuid.NewString() + "@example.com")
}

func SeedAccount(pool *pgxpool.Pool, email string, displayName, password string) (uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, err
	}
	var id uuid.UUID
	err = pool.QueryRow(
		context.Background(),
		`INSERT INTO account (email, display_name, password_hash, status) VALUES ($1, $2, $3, 'active') RETURNING id`,
		email, displayName, string(hash),
	).Scan(&id)
	return id, err
}

func SeedSession(pool *pgxpool.Pool, accountID uuid.UUID) (uuid.UUID, error) {
	var token uuid.UUID
	err := pool.QueryRow(
		context.Background(),
		`INSERT INTO session (expires_at, account_id) VALUES ($1, $2) RETURNING token`,
		time.Now().Add(24*time.Hour), accountID,
	).Scan(&token)
	return token, err
}
