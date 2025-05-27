-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS discente_metricas_repeticao (
    "id"                    uuid            PRIMARY KEY NOT NULL    DEFAULT gen_random_uuid(),
    "created_by_user_id"    uuid            REFERENCES users(id),
    "trilha_id"             uuid            REFERENCES trilha(id) ON DELETE CASCADE,
    "lista_id"              uuid            REFERENCES lista(id) ON DELETE CASCADE,
    "exercicios_base_id"    uuid            REFERENCES exercicios_base(id) ON DELETE CASCADE,    
    "horario_at"            timestamp       NOT NULL,
    "tipo"                  smallint        NOT NULL
);

comment on column discente_metricas_repeticao.tipo is '1: inicio; 2: tentou; 3: rodou';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS discente_metricas_repeticao;
-- +goose StatementEnd
