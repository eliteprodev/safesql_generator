CREATE TABLE bar (id INT NOT NULL, parent_id INT);

-- name: CTERecursive :many
WITH RECURSIVE cte AS (
        SELECT b.* FROM bar AS b
        WHERE b.id = $1
    UNION ALL
        SELECT b.*
        FROM bar AS b, cte AS c
        WHERE b.parent_id = c.id
) SELECT * FROM cte;
