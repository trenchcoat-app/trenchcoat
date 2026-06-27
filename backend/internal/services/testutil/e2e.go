package testutil

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"trenchcoat/internal/api"
	"trenchcoat/internal/db/dbtest"
	"trenchcoat/internal/handlers"
	"trenchcoat/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5/pgxpool"
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
