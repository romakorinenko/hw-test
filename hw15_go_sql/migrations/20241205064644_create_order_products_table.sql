-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_products
(
    id         BIGINT PRIMARY KEY,
    order_id   BIGINT NOT NULL REFERENCES orders (id) ON DELETE CASCADE ON UPDATE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_products;
-- +goose StatementEnd
