-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE users_sequence start 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop sequence users_sequence;
-- +goose StatementEnd
