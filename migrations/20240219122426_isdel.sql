-- +goose Up
-- +goose StatementBegin
ALTER TABLE short ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE short DROP COLUMN is_deleted;
-- +goose StatementEnd
