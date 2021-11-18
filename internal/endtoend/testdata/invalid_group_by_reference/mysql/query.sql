CREATE TABLE authors (
  id   BIGINT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name text    NOT NULL,
  bio  text,
  UNIQUE(name)
);

-- name: ListAuthors :many
SELECT   *
FROM     authors
GROUP BY invalid_reference;
