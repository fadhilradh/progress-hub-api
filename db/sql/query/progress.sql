-- name: CreateProgress :exec
INSERT INTO progress(
    chart_id, 
    range_value, 
    progress_value, 
    created_at,
    progress_no
)
VALUES(
    $1, 
    $2, 
    $3,
    now(),
    $4
) RETURNING *;

-- name: GetProgressByChartID :many
SELECT * FROM progress WHERE chart_id = $1;