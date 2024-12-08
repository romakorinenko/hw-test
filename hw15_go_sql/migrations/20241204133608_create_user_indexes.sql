-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS users_name_idx ON users USING btree (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX users_name_idx;
-- +goose StatementEnd
