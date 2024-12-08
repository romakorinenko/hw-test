-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS order_products_order_id_idx ON order_products USING btree (order_id);

CREATE INDEX IF NOT EXISTS order_products_product_id_idx ON order_products USING btree (product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX order_products_order_id_idx;

DROP INDEX order_products_product_id_idx;
-- +goose StatementEnd
