-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS galleries (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    title TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Downz
-- +goose StatementBegin
DROP TABLE galleries;
-- +goose StatementEnd
