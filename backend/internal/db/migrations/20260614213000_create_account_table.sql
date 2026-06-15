-- +goose Up
CREATE TYPE users_status AS ENUM('active', 'disabled');

CREATE TABLE IF NOT EXISTS account (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  email TEXT NOT NULL UNIQUE,
  display_name TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  status users_status NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS account;
DROP TYPE IF EXISTS users_status;
