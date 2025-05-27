-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS users ADD COLUMN repeticao_espacada_minutos bigint NOT NULL DEFAULT 60;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS users DROP COLUMN repeticao_espacada_minutos;
-- +goose StatementEnd
