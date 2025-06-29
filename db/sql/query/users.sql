-- name: CreateUser :one
INSERT INTO users(id, username, email, password, created_at, updated_at, role)
VALUES($1, $2, $3, $4, $5, $6, $7)
RETURNING id, username, email, role, created_at, updated_at;

-- name: GetUserByUsername :one 
SELECT * FROM users WHERE username = $1;