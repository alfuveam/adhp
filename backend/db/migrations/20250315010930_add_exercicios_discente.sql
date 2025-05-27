-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS exercicios_discente (
    "id"                    uuid            PRIMARY KEY NOT NULL    DEFAULT gen_random_uuid(),
    "created_by_user_id"    uuid            REFERENCES users(id),
    "lista_id"              uuid            REFERENCES lista(id) ON DELETE CASCADE,
    "exercicios_base_id"    uuid            REFERENCES exercicios_base(id) ON DELETE CASCADE,
    "habilitado"            boolean         NOT NULL,
    "codigo_rodou"          boolean         NOT NULL,
    "codigo_base"           text            NOT NULL,
    "create_at"             timestamp       NOT NULL,
    "update_at"             timestamp       NOT NULL,
    UNIQUE (created_by_user_id, exercicios_base_id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exercicios_discente;
-- +goose StatementEnd
