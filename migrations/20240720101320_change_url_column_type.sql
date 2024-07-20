-- +goose Up
-- +goose StatementBegin
ALTER TABLE urls
ALTER COLUMN url 
TYPE TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE urls
ALTER COLUMN url 
TYPE VARCHAR(50);
-- +goose StatementEnd
