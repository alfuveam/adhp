-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS exercicios_repeticao (
    "id"                    uuid            PRIMARY KEY NOT NULL    DEFAULT gen_random_uuid(),
    "created_by_user_id"    uuid            REFERENCES users(id),
    "exercicios_base_id"    uuid            REFERENCES exercicios_base(id) ON DELETE CASCADE,
    "repeticao"             integer         NOT NULL DEFAULT 0,
    "proxima_repeticao"     timestamp       NOT NULL,
    "create_at"             timestamp       NOT NULL,
    "update_at"             timestamp       NOT NULL,
    UNIQUE (created_by_user_id, exercicios_base_id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exercicios_repeticao;
-- +goose StatementEnd
