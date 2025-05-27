-- name: CreateLista :one
INSERT INTO lista (created_by_user_id, trilha_id, order_index, name, description, create_at, update_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAllListas :many
SELECT * FROM lista
ORDER BY id;

-- name: GetListaById :one
SELECT * FROM lista WHERE id = $1;

-- name: GetListaByTrilhaId :many
SELECT * FROM lista WHERE trilha_id = $1 ORDER BY order_index ASC;

-- name: GetListaByOrderIndexAndTrilhaId :one
SELECT * FROM lista WHERE order_index = $1 AND trilha_id = $2;

-- name: GetListaCountByTrilhaId :one
SELECT COUNT(*)::smallint FROM lista WHERE trilha_id = $1;

-- name: UpdateLista :one
UPDATE lista
SET created_by_user_id = $2, trilha_id = $3, order_index = $4, name = $5, description = $6, update_at = $7
WHERE id = $1
RETURNING *;

-- name: DeleteLista :execresult
DELETE FROM lista WHERE id = $1;

