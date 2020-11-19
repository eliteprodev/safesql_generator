-- Comparison Functions and Operators
-- https://www.postgresql.org/docs/current/functions-comparison.html

CREATE TABLE bar (id serial not null);

-- name: GreaterThan :many
SELECT count(*) > 0 FROM bar;

-- name: LessThan :many
SELECT count(*) < 0 FROM bar;

-- name: GreaterThanOrEqual :many
SELECT count(*) >= 0 FROM bar;

-- name: LessThanOrEqual :many
SELECT count(*) <= 0 FROM bar;

-- name: NotEqual :many
SELECT count(*) != 0 FROM bar;

-- name: AlsoNotEqual :many
SELECT count(*) <> 0 FROM bar;

-- name: Equal :many
SELECT count(*) = 0 FROM bar;






