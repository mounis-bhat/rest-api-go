-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tokens (
    hash BYTEA PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expiry TIMESTAMP WITH TIME ZONE NOT NULL,
    scope TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tokens;
-- +goose StatementEnd