CREATE TABLE foo (bar_id serial not null, baz_id serial not null);
CREATE TABLE bar (id serial not null);
CREATE TABLE baz (id serial not null);

-- name: TwoJoins :many
SELECT foo.*
FROM foo
JOIN bar ON bar.id = bar_id
JOIN baz ON baz.id = baz_id;
