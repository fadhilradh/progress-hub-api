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

-- name: GetChartProgressByUserId :many
SELECT c.id AS chart_id, p.id as progress_id, c.range_type, p.range_value, c.progress_name, p.progress_value,  
p.updated_at AS progress_updated_at
FROM charts c
INNER JOIN progress p ON c.id = p.chart_id 
WHERE c.user_id = $1;