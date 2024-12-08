-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products
(
    id    BIGINT PRIMARY KEY,
    name  varchar(255)   NOT NULL UNIQUE,
    price NUMERIC(10, 2) NOT NULL CHECK (price > 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
