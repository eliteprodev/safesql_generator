CREATE TABLE foo (num integer not null);

-- name: Math :many
SELECT *, num / 1024 as division FROM foo;
