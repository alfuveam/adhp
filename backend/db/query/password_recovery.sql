-- name: UpsertPasswordRecovery :exec
    INSERT INTO password_recovery (user_id, code_hash, attempts, createat, expires_at)
    VALUES ($1, $2, 0, $3, $4)
    ON CONFLICT (user_id) DO UPDATE
    SET code_hash = EXCLUDED.code_hash, attempts = 0, createat = EXCLUDED.createat, expires_at = EXCLUDED.expires_at;

-- name: GetPasswordRecovery :one
    SELECT user_id, code_hash, attempts, createat, expires_at FROM password_recovery WHERE user_id = $1;

-- name: IncrementPasswordRecoveryAttempts :exec
    UPDATE password_recovery SET attempts = attempts + 1 WHERE user_id = $1;

-- name: DeletePasswordRecovery :exec
    DELETE FROM password_recovery WHERE user_id = $1;
