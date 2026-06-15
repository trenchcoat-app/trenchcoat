package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"divvy/config"
)

func OpenDB() *pgxpool.Pool {
	dsn := DSN()

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	return pool
}

func DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.AppConfig.POSTGRES_USER,
		config.AppConfig.POSTGRES_PASSWORD,
		config.AppConfig.POSTGRES_HOST,
		config.AppConfig.POSTGRES_PORT,
		config.AppConfig.POSTGRES_DB,
	)
}
