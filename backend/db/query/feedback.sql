-- name: CreateFeedback :one
INSERT INTO feedback (created_by_user_id, exercicios_base_id, descricao, create_at, update_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAllFeedbacks :many
SELECT * FROM feedback
ORDER BY id;

-- name: GetFeedbackById :one
SELECT * FROM feedback WHERE id = $1;

-- name: GetFeedbackByExerciciosBaseId :many
SELECT * FROM feedback WHERE exercicios_base_id = $1;

-- name: GetRandomFeedbackByExerciciosBaseId :one
SELECT * FROM feedback WHERE exercicios_base_id = $1 ORDER BY random() LIMIT 1;

-- name: UpdateFeedbackDescricaoById :one
UPDATE feedback
SET descricao = $2, update_at = $3
WHERE id = $1
RETURNING *;

-- name: UpdateFeedback :one
UPDATE feedback
SET created_by_user_id = $2, exercicios_base_id = $3, descricao = $4, update_at = $5
WHERE id = $1
RETURNING *;

-- name: DeleteFeedback :execresult
DELETE FROM feedback WHERE id = $1;