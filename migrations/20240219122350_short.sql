-- +goose Up
-- +goose StatementBegin
CREATE TABLE short (id SERIAL , brief TEXT NOT NULL, origin TEXT NOT NULL UNIQUE);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE short;
-- +goose StatementEnd
