CREATE TABLE foo (a text, b integer);

-- name: InsertValues :exec
INSERT INTO public.foo (a, b) VALUES ($1, $2);
