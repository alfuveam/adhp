-- name: CreateTrilha :one
    INSERT INTO trilha (created_by_user_id, name, description, tipo_da_linguagem, create_at, update_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *;

-- name: GetAllTrilhas :many
    SELECT * FROM trilha
    ORDER BY id;

-- name: GetTrilhaById :one
    SELECT * FROM trilha WHERE id = $1;

-- name: UpdateTrilha :one
    UPDATE trilha
    SET created_by_user_id = $2, name = $3, description = $4, tipo_da_linguagem = $5, update_at = $6
    WHERE id = $1
    RETURNING *;

-- name: DeleteTrilha :execresult
    DELETE FROM trilha WHERE id = $1;

