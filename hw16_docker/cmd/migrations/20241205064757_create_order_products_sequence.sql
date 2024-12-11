-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE order_products_sequence start 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop sequence order_products_sequence;
-- +goose StatementEnd
