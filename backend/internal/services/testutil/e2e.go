package testutil

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"trenchcoat/internal/api"
	"trenchcoat/internal/db/dbtest"
	"trenchcoat/internal/handlers"
	"trenchcoat/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetE2EPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	return dbtest.SetupDB(t)
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
