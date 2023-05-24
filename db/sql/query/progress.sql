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

-- name: EditProgressByID :exec
UPDATE progress SET range_value = $2, progress_value = COALESCE($3, progress_value), updated_at = now() WHERE id = $1 RETURNING *;