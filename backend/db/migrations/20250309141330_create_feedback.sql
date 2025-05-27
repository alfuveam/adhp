-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS feedback (
    "id"                    uuid            PRIMARY KEY NOT NULL    DEFAULT gen_random_uuid(),
    "created_by_user_id"    uuid            REFERENCES users(id),
    "exercicios_base_id"    uuid            REFERENCES exercicios_base(id) ON DELETE CASCADE,
    "descricao"             text            NOT NULL,
    "create_at"             timestamp       NOT NULL,
    "update_at"             timestamp       NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feedback;
-- +goose StatementEnd
