CREATE TABLE foo (a text, b text);
CREATE TABLE bar (c text, d text);
-- name: StarExpansionJoin :many
SELECT * FROM foo, bar;
