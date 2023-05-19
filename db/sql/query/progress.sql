-- name: CreateProgress :exec
INSERT INTO progress(
    chart_id, 
    range_value, 
    progress_value, 
    created_at
)
VALUES(
    $1, 
    $2, 
    $3,
    now()
) RETURNING *;