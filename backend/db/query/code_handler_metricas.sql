-- name: CreateCodeHandlerMetricas :one
INSERT INTO code_handler_metricas (created_by_user_id, trilha_id, lista_id, exercicios_base_id, horario_at, tipo, linguagem)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAllCodeHandlerMetricas :many
SELECT * FROM code_handler_metricas
ORDER BY id;

-- name: GetCodeHandlerMetricasById :one
SELECT * FROM code_handler_metricas WHERE id = $1;

-- name: UpdateCodeHandlerMetricas :one
UPDATE code_handler_metricas
SET created_by_user_id = $2, trilha_id = $3, lista_id = $4, exercicios_base_id = $5, horario_at = $6, tipo = $7, linguagem = $8
WHERE id = $1
RETURNING *;

-- name: DeleteCodeHandlerMetricas :execresult
DELETE FROM code_handler_metricas WHERE id = $1;

