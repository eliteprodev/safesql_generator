CREATE SCHEMA myschema;
CREATE TABLE myschema.foo (a text, b integer);

-- name: InsertValues :copyfrom
-- InsertValues inserts multiple values using copy.
INSERT INTO myschema.foo (a, b) VALUES ($1, $2);

-- name: InsertSingleValue :copyfrom
-- InsertSingleValue inserts a single value using copy.
INSERT INTO myschema.foo (a) VALUES ($1);
