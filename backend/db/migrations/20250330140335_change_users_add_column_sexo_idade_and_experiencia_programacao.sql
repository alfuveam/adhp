-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS users ADD COLUMN idade smallint NOT NULL DEFAULT 0;
ALTER TABLE IF EXISTS users ADD COLUMN experiencia_programacao boolean NOT NULL DEFAULT false;
ALTER TABLE IF EXISTS users ADD COLUMN sexo smallint NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS users DROP COLUMN idade;
ALTER TABLE IF EXISTS users DROP COLUMN experiencia_programacao;
ALTER TABLE IF EXISTS users DROP COLUMN sexo;
-- +goose StatementEnd
