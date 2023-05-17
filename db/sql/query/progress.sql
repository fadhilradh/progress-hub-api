-- name: CreateProgress :one
INSERT INTO progress(
    user_id, 
    range_type, 
    range_value, 
    progress_name, 
    progress_value, 
    created_at
)
VALUES(
    $1, 
    $2, 
    $3,
    $4,
    $5,
    now()
) RETURNING *;