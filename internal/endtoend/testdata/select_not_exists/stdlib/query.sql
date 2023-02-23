CREATE TABLE bar (id serial not null);

-- name: BarNotExists :one
SELECT
    NOT EXISTS (
        SELECT
            1
        FROM
            bar
        where
            id = $1
    );
