-- name: CreateUser :one
    INSERT INTO users (completename, email, password, createat, lastlogin, usertype, repeticao_espacada_minutos) 
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING *;
        
-- name: OnLoginCheckEmailAndPasswordId :one
    SELECT id FROM users WHERE email = $1 AND password = $2;
    
-- name: OnLoginCheckEmailAndPasswordUser :one
    SELECT id, completename, email, password, createat, lastlogin, usertype, repeticao_espacada_minutos FROM users WHERE email = $1;

-- name: GetUserById :one
    SELECT id, completename, email, password, createat, lastlogin, usertype, repeticao_espacada_minutos FROM users WHERE id = $1;

-- name: UpdateUserRepeticao :one
    UPDATE users
    SET repeticao_espacada_minutos = $2
    WHERE id = $1
    RETURNING *;
