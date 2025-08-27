-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE apps (
    id TEXT UNIQUE PRIMARY KEY,
    name TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO apps (id, name, thumbnail_url) VALUES ('1458862350', 'Hevy - Workout Tracker Gym Log', 'https://is1-ssl.mzstatic.com/image/thumb/Purple211/v4/ab/9b/08/ab9b08ae-6dd2-c978-899b-e0edcf1aba9c/AppIcon-0-0-1x_U007emarketing-0-7-0-sRGB-85-220.png/512x512bb.jpg');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE apps;
-- +goose StatementEnd
