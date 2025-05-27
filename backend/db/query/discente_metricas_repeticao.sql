-- name: CreateDiscenteMetricasRepeticao :one
INSERT INTO discente_metricas_repeticao (created_by_user_id, trilha_id, lista_id, exercicios_base_id, horario_at, tipo, repeticao_espacada_minutos)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetDiscenteMetricasRepeticaoById :one
SELECT * FROM discente_metricas_repeticao WHERE id = $1;

-- name: GetAllDiscenteMetricasRepeticao :many
SELECT * FROM discente_metricas_repeticao
ORDER BY id;

-- name: UpdateDiscenteMetricasRepeticao :one
UPDATE discente_metricas_repeticao
SET created_by_user_id = $2, trilha_id = $3, lista_id = $4, exercicios_base_id = $5, horario_at = $6, tipo = $7, repeticao_espacada_minutos = $8
WHERE id = $1
RETURNING *;

-- name: DeleteDiscenteMetricasRepeticao :execresult
DELETE FROM discente_metricas_repeticao WHERE id = $1;

