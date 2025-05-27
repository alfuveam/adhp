-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS discente_metricas_feedback (
    "id"                    uuid            PRIMARY KEY NOT NULL    DEFAULT gen_random_uuid(),
    "created_by_user_id"    uuid            REFERENCES users(id),
    "trilha_id"             uuid            REFERENCES trilha(id) ON DELETE CASCADE,
    "lista_id"              uuid            REFERENCES lista(id) ON DELETE CASCADE,
    "exercicios_base_id"    uuid            REFERENCES exercicios_base(id) ON DELETE CASCADE,    
    "horario_at"            timestamp       NOT NULL,
    "tipo_exercicio"        smallint        NOT NULL
);

comment on column discente_metricas_repeticao.tipo is '1: do exercício; 2: da repetição';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS discente_metricas_feedback;
-- +goose StatementEnd
