-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS products_name_idx ON products USING btree (name);

CREATE INDEX IF NOT EXISTS products_price_idx ON products USING btree (price);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX products_price_idx;

DROP INDEX products_name_idx;
-- +goose StatementEnd
