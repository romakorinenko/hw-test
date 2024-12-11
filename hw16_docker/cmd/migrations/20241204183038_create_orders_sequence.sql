-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE orders_sequence start 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop sequence orders_sequence;
-- +goose StatementEnd
