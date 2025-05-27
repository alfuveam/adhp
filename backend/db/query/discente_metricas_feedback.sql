-- name: CreateDiscenteMetricasFeedback :one
INSERT INTO discente_metricas_feedback (created_by_user_id, trilha_id, lista_id, exercicios_base_id, horario_at, tipo_exercicio)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllDiscenteMetricasFeedback :many
SELECT * FROM discente_metricas_feedback
ORDER BY id;

-- name: GetDiscenteMetricasFeedbackById :one
SELECT * FROM discente_metricas_feedback WHERE id = $1;

-- name: GetDiscenteMetricasFeedbackByTrilhaId :many
SELECT * FROM discente_metricas_feedback WHERE trilha_id = $1;

-- name: GetDiscenteMetricasFeedbackByListaId :many
SELECT * FROM discente_metricas_feedback WHERE lista_id = $1;

-- name: GetDiscenteMetricasFeedbackByExerciciosBaseId :many
SELECT * FROM discente_metricas_feedback WHERE exercicios_base_id = $1;

-- name: UpdateDiscenteMetricasFeedback :one
UPDATE discente_metricas_feedback
SET created_by_user_id = $2, trilha_id = $3, lista_id = $4, exercicios_base_id = $5, horario_at = $6, tipo_exercicio = $7
WHERE id = $1
RETURNING *;

-- name: DeleteDiscenteMetricasFeedback :execresult
DELETE FROM discente_metricas_feedback WHERE id = $1;