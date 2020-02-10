CREATE TABLE foo (a text, b text);
CREATE TABLE bar (c text, d text);

-- name: StarExpansionCTE :many
WITH cte AS (SELECT * FROM foo) SELECT * FROM bar;

-- name: StarExpansionTwoCTE :many
WITH 
  a AS (SELECT * FROM foo),
  b AS (SELECT 1::int as bar, * FROM a)
SELECT * FROM b;
