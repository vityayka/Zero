-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS oauth_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    provider TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE oauth_tokens;
-- +goose StatementEnd
