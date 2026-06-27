package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedAccount(t *testing.T, pool *pgxpool.Pool, email, displayName, password string) uuid.UUID {
	t.Helper()
	hash := HashPassword(t, password)
	var id uuid.UUID
	err := pool.QueryRow(
		context.Background(),
		`INSERT INTO account (email, display_name, password_hash, status) VALUES ($1, $2, $3, 'active') RETURNING id`,
		email, displayName, hash,
	).Scan(&id)
	require.NoError(t, err)
	return id
}

func SeedSession(t *testing.T, pool *pgxpool.Pool, accountID uuid.UUID) uuid.UUID {
	t.Helper()
	var token uuid.UUID
	err := pool.QueryRow(
		context.Background(),
		`INSERT INTO session (expires_at, account_id) VALUES ($1, $2) RETURNING token`,
		time.Now().Add(24*time.Hour), accountID,
	).Scan(&token)
	require.NoError(t, err)
	return token
}
