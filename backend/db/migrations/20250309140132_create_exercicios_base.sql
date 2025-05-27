-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS exercicios_base (
    "id"                    uuid            PRIMARY KEY NOT NULL    DEFAULT gen_random_uuid(),
    "created_by_user_id"    uuid            REFERENCES users(id),
    "lista_id"              uuid            REFERENCES lista(id) ON DELETE CASCADE,
    "order_index"           smallint        NOT NULL,
    "titulo"                text            NOT NULL,
    "codigo_base"           text            NOT NULL,
    "codigo_teste"          text            NOT NULL,
    "create_at"             timestamp       NOT NULL,
    "update_at"             timestamp       NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exercicios_base;
-- +goose StatementEnd
