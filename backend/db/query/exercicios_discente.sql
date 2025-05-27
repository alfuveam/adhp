-- name: CreateExerciciosDiscente :one
INSERT INTO exercicios_discente (created_by_user_id, lista_id, exercicios_base_id, codigo_base, codigo_rodou, create_at, update_at, habilitado)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAllExerciciosDiscente :many
SELECT * FROM exercicios_discente WHERE created_by_user_id = $1
ORDER BY id;

-- name: GetExerciciosDiscenteById :one
SELECT * FROM exercicios_discente WHERE id = $1 AND created_by_user_id = $2;

-- name: GetExerciciosDiscenteByExerciciosBaseId :one
SELECT * FROM exercicios_discente WHERE exercicios_base_id = $1 AND created_by_user_id = $2;

-- name: CheckIfHasExerciciosDiscenteByExerciciosBaseId :one
SELECT EXISTS (SELECT * FROM exercicios_discente WHERE exercicios_base_id = $1 AND created_by_user_id = $2) as exercicio_existe;

-- name: GetExerciciosDiscenteByListaId :many
SELECT * FROM exercicios_discente WHERE lista_id = $1 AND created_by_user_id = $2;

-- name: UpsertExerciciosDiscente :one
INSERT INTO exercicios_discente (created_by_user_id, lista_id, exercicios_base_id, codigo_base, codigo_rodou, create_at, update_at, habilitado)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (created_by_user_id, exercicios_base_id) DO UPDATE
SET lista_id = excluded.lista_id, codigo_base = excluded.codigo_base, codigo_rodou = excluded.codigo_rodou, update_at = excluded.update_at, habilitado = excluded.habilitado
RETURNING *, (CASE WHEN xmax = 0 THEN 'insert' ELSE 'update' END) AS tipo_operacao;

-- name: UpdateExerciciosDiscente :one
UPDATE exercicios_discente
SET created_by_user_id = $2, lista_id = $3, exercicios_base_id = $4, codigo_base = $5, codigo_rodou = $6, update_at = $7, habilitado = $8
WHERE id = $1
RETURNING *;

-- name: DeleteExerciciosDiscente :execresult
DELETE FROM exercicios_discente WHERE id = $1;