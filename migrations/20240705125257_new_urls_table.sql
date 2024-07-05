-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    alias VARCHAR(10) UNIQUE NOT NULL,
    url VARCHAR(50) NOT NULL,
    expire TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd
