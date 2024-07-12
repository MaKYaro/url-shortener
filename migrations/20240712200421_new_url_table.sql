-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
    alias VARCHAR(10) PRIMARY KEY,
    url VARCHAR(50) NOT NULL,
    expire TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd
