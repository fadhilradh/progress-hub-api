-- name: CreateChart :one
INSERT INTO charts(
    user_id, 
    range_type, 
    progress_name, 
    created_at
)
VALUES(
    $1, 
    $2, 
    $3,
    now()
) RETURNING *;