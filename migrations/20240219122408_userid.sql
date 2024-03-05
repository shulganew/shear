-- +goose Up
-- +goose StatementBegin
ALTER TABLE short ADD COLUMN  user_id TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE short DROP COLUMN user_id;
-- +goose StatementEnd
