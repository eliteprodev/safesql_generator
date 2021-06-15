CREATE TABLE bar (id serial not null);

-- name: BarExists :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            bar
        where
            id = $1
    );
