package dbtest

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	Pool    *pgxpool.Pool
	cleanup func()
}

func SetupDB(t testing.TB) *pgxpool.Pool {
	t.Helper()

	db, err := startContainer(t)
	require.NoError(t, err, "Failed to start postgres container")

	t.Cleanup(db.cleanup)

	return db.Pool
}

func startContainer(t testing.TB) (*TestDB, error) {
	ctx := context.Background()

	container, err := postgres.Run(ctx, "postgres:18.4",
		postgres.WithDatabase("trenchcoat_test"),
		postgres.WithUsername("trenchcoat"),
		postgres.WithPassword("trenchcoat"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err, "Failed to run postgres container")

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "Failed to retrieve postgres connection string")

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err, "Failed to create new postgres pool")

	sqlDB, err := sql.Open("pgx", connStr)
	require.NoError(t, err, "Failed to open postgres connection via connection string")

	err = goose.SetDialect("postgres")
	require.NoError(t, err, "Failed to set Goose dialect")

	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(filename), "..", "migrations")

	err = goose.Up(sqlDB, migrationsDir)
	require.NoError(t, err, "Failed to apply migrations via Goose")

	// We do the following inline instead of extracting the function to get access to pool and sqlDB without params
	cleanup := func() {
		pool.Close()
		sqlDB.Close()
		if err := container.Terminate(ctx); err != nil {
			panic("failed to terminate container: " + err.Error())
		}
	}

	return &TestDB{Pool: pool, cleanup: cleanup}, nil
}
