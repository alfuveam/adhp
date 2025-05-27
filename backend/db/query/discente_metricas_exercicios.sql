-- name: CreateDiscenteMetricasExercicios :one
INSERT INTO discente_metricas_exercicios (created_by_user_id, trilha_id, lista_id, exercicios_base_id, horario_at, tipo)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllDiscenteMetricasExercicios :many
SELECT * FROM discente_metricas_exercicios
ORDER BY id;

-- name: GetDiscenteMetricasExerciciosById :one
SELECT * FROM discente_metricas_exercicios WHERE id = $1;

-- name: UpdateDiscenteMetricasExercicios :one
UPDATE discente_metricas_exercicios
SET created_by_user_id = $2, trilha_id = $3, lista_id = $4, exercicios_base_id = $5, horario_at = $6, tipo = $7
WHERE id = $1
RETURNING *;

-- name: DeleteDiscenteMetricasExercicios :execresult
DELETE FROM discente_metricas_exercicios WHERE id = $1;
