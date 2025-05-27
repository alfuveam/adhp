-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS code_handler_metricas (
    "id"                    uuid            PRIMARY KEY NOT NULL    DEFAULT gen_random_uuid(),
    "created_by_user_id"    uuid            REFERENCES users(id),
    "trilha_id"             uuid            REFERENCES trilha(id) ON DELETE CASCADE,
    "lista_id"              uuid            REFERENCES lista(id) ON DELETE CASCADE,
    "exercicios_base_id"    uuid            REFERENCES exercicios_base(id) ON DELETE CASCADE,    
    "horario_at"            timestamp       NOT NULL,
    "tipo"                  smallint        NOT NULL,
    "linguagem"             smallint        NOT NULL
);

comment on column code_handler_metricas.tipo is '1: inicio; 2: falhou; 3: rodou';
comment on column code_handler_metricas.linguagem is '1: golang; 2: python';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS code_handler_metricas;
-- +goose StatementEnd
