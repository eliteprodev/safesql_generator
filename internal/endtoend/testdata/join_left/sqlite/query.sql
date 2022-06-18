-- https://github.com/kyleconroy/sqlc/issues/604
CREATE TABLE users (
  user_id    INT PRIMARY KEY,
  city_id    INT -- nullable
);
CREATE TABLE cities (
  city_id    INT PRIMARY KEY,
  mayor_id   INT NOT NULL
);
CREATE TABLE mayors (
  mayor_id   INT PRIMARY KEY,
  full_name  TEXT NOT NULL
);

-- name: GetMayors :many
SELECT
    user_id,
    mayors.full_name
FROM users
LEFT JOIN cities USING (city_id)
INNER JOIN mayors USING (mayor_id);

-- name: GetMayorsOptional :many
SELECT
    user_id,
    cities.city_id,
    mayors.full_name
FROM users
LEFT JOIN cities USING (city_id)
LEFT JOIN mayors USING (mayor_id);

-- https://github.com/kyleconroy/sqlc/issues/1334
CREATE TABLE authors (
  id        INT PRIMARY KEY,
  name      TEXT NOT NULL,
  parent_id INT -- nullable
);

CREATE TABLE super_authors (
  super_id        INT PRIMARY KEY,
  super_name      TEXT NOT NULL,
  super_parent_id INT -- nullable
);

-- name: AllAuthors :many
SELECT  *
FROM    authors AS a
        LEFT JOIN authors AS p
            ON a.parent_id = p.id;

-- name: AllAuthorsAliases :many
SELECT  *
FROM    authors AS a
        LEFT JOIN authors AS p
            ON a.parent_id = p.id;

-- name: AllSuperAuthors :many
SELECT  *
FROM    authors
        LEFT JOIN super_authors
            ON authors.parent_id = super_authors.super_id;

-- name: AllAuthorsAliases2 :many
SELECT  a.*, p.*
FROM    authors AS a
        LEFT JOIN authors AS p
            ON a.parent_id = p.id;

-- name: AllSuperAuthorsAliases :many
SELECT  *
FROM    authors AS a
        LEFT JOIN super_authors AS sa
            ON a.parent_id = sa.super_id;

-- name: AllSuperAuthorsAliases2 :many
SELECT  a.*, sa.*
FROM    authors AS a
        LEFT JOIN super_authors AS sa
            ON a.parent_id = sa.super_id;

-- https://github.com/kyleconroy/sqlc/issues/1334
CREATE TABLE users_2 (
    user_id           INT PRIMARY KEY,
    user_nickname     VARCHAR(30) UNIQUE NOT NULL,
    user_email        TEXT UNIQUE        NOT NULL,
    user_display_name TEXT               NOT NULL,
    user_password     TEXT               ,
    user_google_id    TEXT UNIQUE        ,
    user_apple_id     TEXT UNIQUE        ,
    user_bio          VARCHAR(160)       NOT NULL DEFAULT '',
    user_created_at   TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_avatar_id    INT UNIQUE        
);

CREATE TABLE media (
    media_id         INT PRIMARY KEY,
    media_created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    media_hash       TEXT      NOT NULL,
    media_directory  TEXT      NOT NULL,
    media_author_id  INT       NOT NULL,
    media_width      INT       NOT NULL,
    media_height     INT       NOT NULL
);

-- name: GetSuggestedUsersByID :many
SELECT  DISTINCT u.*, m.*
FROM    users_2 AS u
        LEFT JOIN media AS m
            ON u.user_avatar_id = m.media_id
WHERE   u.user_id != @user_id;
