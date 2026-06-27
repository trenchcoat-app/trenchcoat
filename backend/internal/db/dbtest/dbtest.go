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

	db, err := startContainer()
	require.NoError(t, err, "Failed to start postgres container")

	t.Cleanup(db.cleanup)

	return db.Pool
}

func SetupDBMain() (*pgxpool.Pool, func(), error) {
	db, err := startContainer()
	if err != nil {
		return nil, nil, err
	}
	return db.Pool, db.cleanup, nil
}

func startContainer() (*TestDB, error) {
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
	if err != nil {
		return nil, err
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		return nil, err
	}

	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(filename), "..", "migrations")

	err = goose.Up(sqlDB, migrationsDir)
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		pool.Close()
		sqlDB.Close()
		if err := container.Terminate(ctx); err != nil {
			panic("failed to terminate container: " + err.Error())
		}
	}

	return &TestDB{Pool: pool, cleanup: cleanup}, nil
}
