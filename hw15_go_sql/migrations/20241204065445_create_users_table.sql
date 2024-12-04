-- +goose Up
-- +goose StatementBegin
CREATE table IF NOT EXISTS users
(
    id       BIGINT PRIMARY KEY,
    name     VARCHAR(255) NOT NULL DEFAULT '',
    email    VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
