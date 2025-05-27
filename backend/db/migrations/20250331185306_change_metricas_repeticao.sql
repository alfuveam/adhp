-- +goose Up
-- +goose StatementBegin
ALTER TABLE discente_metricas_repeticao ADD COLUMN repeticao_espacada_minutos bigint NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE discente_metricas_repeticao DROP COLUMN repeticao_espacada_minutos
-- +goose StatementEnd
