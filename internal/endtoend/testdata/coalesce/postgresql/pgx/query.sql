CREATE TABLE foo (bar text, bat text not null);

-- name: Coalesce :many
SELECT coalesce(bar, '') as login
FROM foo;

-- name: CoalesceColumns :many
SELECT bar, bat, coalesce(bar, bat)
FROM foo;
