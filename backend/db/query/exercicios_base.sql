-- name: CreateExerciciosBase :one
INSERT INTO exercicios_base (created_by_user_id, lista_id, order_index, titulo, codigo_base, codigo_teste, create_at, update_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAllExerciciosBase :many
SELECT * FROM exercicios_base
ORDER BY id;

-- name: GetExerciciosBaseById :one
SELECT * FROM exercicios_base WHERE id = $1;

-- name: GetExerciciosBaseByListaId :many
SELECT * FROM exercicios_base WHERE lista_id = $1 ORDER BY order_index ASC;

-- name: GetExercicioByOrderIndexAndListaId :one
SELECT * FROM exercicios_base WHERE order_index = $1 AND lista_id = $2;

-- name: GetExerciciosBaseCountByListaId :one
SELECT COUNT(*)::smallint FROM exercicios_base WHERE lista_id = $1;

-- name: UpdateExerciciosBase :one
UPDATE exercicios_base
SET created_by_user_id = $2, lista_id = $3, order_index = $4, titulo = $5, codigo_base = $6, codigo_teste = $7, update_at = $8
WHERE id = $1
RETURNING *;

-- name: DeleteExerciciosBase :execresult
DELETE FROM exercicios_base WHERE id = $1;

