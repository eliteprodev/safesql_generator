CREATE TABLE bar (id serial primary key not null, ready bool not null);

-- name: DeleteReadyWithCTE :many
WITH ready_ids AS (
	SELECT id FROM bar WHERE ready
)
DELETE FROM bar WHERE id IN (SELECT * FROM ready_ids)
RETURNING id;
