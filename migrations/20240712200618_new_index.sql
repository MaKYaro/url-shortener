-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS alias ON urls(alias);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS alias;
-- +goose StatementEnd
