-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS orders_user_id_idx ON orders USING btree (user_id);

CREATE INDEX IF NOT EXISTS orders_order_date_idx ON orders USING btree (order_date);

CREATE INDEX IF NOT EXISTS orders_total_amount_idx ON orders USING btree (total_amount);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX orders_user_id_idx;

DROP INDEX orders_order_date_idx;

DROP INDEX orders_total_amount_idx;
-- +goose StatementEnd
