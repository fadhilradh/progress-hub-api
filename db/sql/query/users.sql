-- name: CreateUser :one
INSERT INTO users(username, email, password, created_at, updated_at, role)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING id, username, email, role, created_at, updated_at;