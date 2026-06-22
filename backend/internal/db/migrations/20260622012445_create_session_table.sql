-- +goose Up
CREATE TABLE session (
    token UUID PRIMARY KEY DEFAULT uuidv7(),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address TEXT,
    user_agent TEXT,
    account_id UUID NOT NULL REFERENCES account(id) ON DELETE CASCADE
);

CREATE INDEX session_expires_at_idx ON session (expires_at);
CREATE INDEX session_account_id_idx ON session (account_id);

-- +goose Down
DROP TABLE IF EXISTS session;
