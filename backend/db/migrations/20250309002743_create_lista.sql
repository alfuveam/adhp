-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS lista (
    "id"                    uuid            PRIMARY KEY NOT NULL    DEFAULT gen_random_uuid(),
    "created_by_user_id"    uuid            REFERENCES users(id),
    "trilha_id"             uuid            REFERENCES trilha(id) ON DELETE CASCADE,
    "order_index"           smallint        NOT NULL,
    "name"                  varchar(255)    NOT NULL,
    "description"           text,
    "create_at"             timestamp       NOT NULL,
    "update_at"             timestamp       NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lista;
-- +goose StatementEnd
