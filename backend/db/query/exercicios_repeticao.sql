-- name: CreateExercRep :one
    INSERT INTO exercicios_repeticao (created_by_user_id, exercicios_base_id, repeticao, proxima_repeticao, create_at, update_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *;

-- name: GetAllExercRep :many
    SELECT * FROM exercicios_repeticao
    ORDER BY id;

-- name: GetExercRepByIdAndUserId :one
    SELECT * FROM exercicios_repeticao WHERE id = $1 AND created_by_user_id = $2;

-- name: GetExercRepByExercBaseIdAndUserId :many
    SELECT * FROM exercicios_repeticao WHERE exercicios_base_id = $1 AND created_by_user_id = $2 ORDER BY proxima_repeticao ASC;

-- name: GetExercRepByUserId :many
    SELECT * FROM exercicios_repeticao WHERE created_by_user_id = $1 ORDER BY proxima_repeticao ASC;

-- name: UpdateExercRep :one
    UPDATE exercicios_repeticao
    SET created_by_user_id = $2, exercicios_base_id = $3, repeticao = $4, proxima_repeticao = $5, update_at = $6
    WHERE id = $1
    RETURNING *;

-- name: UpdateExercRepRepeticao :one
    UPDATE exercicios_repeticao
    SET repeticao = exercicios_repeticao.repeticao + 1, proxima_repeticao = $3, update_at = $4
    WHERE created_by_user_id = $1 AND exercicios_base_id = $2
    RETURNING *;

-- name: UpsertExercRepNaoSomaRep :one
    INSERT INTO exercicios_repeticao (created_by_user_id, exercicios_base_id, proxima_repeticao, create_at, update_at)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (created_by_user_id, exercicios_base_id) DO UPDATE
    SET proxima_repeticao = excluded.proxima_repeticao, update_at = excluded.update_at
    RETURNING *;

-- name: DeleteExercRep :execresult
    DELETE FROM exercicios_repeticao WHERE id = $1;
