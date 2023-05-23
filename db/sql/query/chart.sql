-- name: CreateChart :one
INSERT INTO charts(
    user_id, 
    range_type, 
    progress_name, 
    created_at,
    colors
)
VALUES(
    $1, 
    $2, 
    $3,
    now(),
    $4
) RETURNING *;

-- name: GetChartProgressByUserId :many
SELECT c.id as chart_id, c.colors as chart_color, p.id as progress_id, c.range_type, p.range_value, c.progress_name, p.progress_value,  
p.updated_at as progress_updated_at, p.progress_no
FROM charts c
INNER JOIN progress p ON c.id = p.chart_id 
WHERE c.user_id = $1
ORDER BY chart_id DESC
;