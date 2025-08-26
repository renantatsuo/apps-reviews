-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE reviews (
    id TEXT UNIQUE PRIMARY KEY,
    app_id TEXT NOT NULL,
    author TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    rating INTEGER NOT NULL,
    sent_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE reviews;
-- +goose StatementEnd
