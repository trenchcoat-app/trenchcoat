-- +goose Up
CREATE TYPE account_status AS ENUM('active', 'disabled');
CREATE TYPE account_auth_provider AS ENUM('email'); -- extendable to add providers like 'google', 'apple'

CREATE TABLE IF NOT EXISTS account (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  email TEXT NOT NULL UNIQUE,
  verified BOOLEAN,
  status account_status NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS account_auth (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  account_id UUID UNIQUE REFERENCES account (id) ON DELETE CASCADE,
  provider account_auth_provider NOT NULL,
  provider_user_id TEXT, -- null for email account_auth_provider
  password_hash TEXT, -- non-null for email account_auth_provider
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT chk_provider_fields CHECK (
    (provider = 'email' AND password_hash IS NOT NULL AND provider_user_id IS NULL) OR
    (provider != 'email' AND password_hash IS NULL AND provider_user_id IS NOT NULL)
  )
);

CREATE TABLE IF NOT EXISTS account_profile (
  account_id UUID PRIMARY KEY REFERENCES account (id) ON DELETE CASCADE,
  username TEXT UNIQUE,
  avatar_url TEXT,
);

-- +goose Down
DROP TABLE IF EXISTS account_profile;
DROP TABLE IF EXISTS account_auth;
DROP TABLE IF EXISTS account;
DROP TYPE IF EXISTS account_auth_provider;
DROP TYPE IF EXISTS account_status;