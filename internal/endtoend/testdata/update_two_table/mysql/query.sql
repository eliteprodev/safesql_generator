-- https://github.com/kyleconroy/sqlc/issues/1590
CREATE TABLE authors (
  name text NOT NULL,
  deleted_at datetime NOT NULL,
  created_at datetime NOT NULL,
  updated_at datetime NOT NULL
);

CREATE TABLE books (
  is_amazing tinyint(1) NOT NULL,
  deleted_at datetime NOT NULL,
  created_at datetime NOT NULL,
  updated_at datetime NOT NULL
);

-- name: DeleteAuthor :exec
UPDATE
  authors,
  books
SET
  authors.deleted_at = now(),
  books.deleted_at = now()
WHERE
  books.is_amazing = 1
  AND authors.name = sqlc.arg(name);