-- +goose Up
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN cpf;
ALTER TABLE users DROP COLUMN phone;
ALTER TABLE users DROP COLUMN country;
ALTER TABLE users DROP COLUMN idade;
ALTER TABLE users DROP COLUMN sexo;
ALTER TABLE users DROP COLUMN experiencia_programacao;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN cpf                        varchar(11);
ALTER TABLE users ADD COLUMN phone                      varchar(50);
ALTER TABLE users ADD COLUMN country                    smallint        NOT NULL;
ALTER TABLE users ADD COLUMN idade                      smallint NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN sexo                       smallint NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN experiencia_programacao    boolean NOT NULL DEFAULT false;
-- +goose StatementEnd
