-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS trilha ADD COLUMN IF NOT EXISTS tipo_da_linguagem smallint NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS trilha DROP COLUMN IF EXISTS tipo_da_linguagem;
-- +goose StatementEnd
