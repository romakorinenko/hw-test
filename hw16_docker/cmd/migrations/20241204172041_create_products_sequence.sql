-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE products_sequence start 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop sequence products_sequence;
-- +goose StatementEnd
