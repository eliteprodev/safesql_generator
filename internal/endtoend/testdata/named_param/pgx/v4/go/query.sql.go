// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package querytest

import (
	"context"
)

const atParams = `-- name: AtParams :many
SELECT name FROM foo WHERE name = $1 AND $2::bool
`

type AtParamsParams struct {
	Slug   string
	Filter bool
}

func (q *Queries) AtParams(ctx context.Context, arg AtParamsParams) ([]string, error) {
	rows, err := q.db.Query(ctx, atParams, arg.Slug, arg.Filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const funcParams = `-- name: FuncParams :many
SELECT name FROM foo WHERE name = $1 AND $2::bool
`

type FuncParamsParams struct {
	Slug   string
	Filter bool
}

func (q *Queries) FuncParams(ctx context.Context, arg FuncParamsParams) ([]string, error) {
	rows, err := q.db.Query(ctx, funcParams, arg.Slug, arg.Filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertAtParams = `-- name: InsertAtParams :one
INSERT INTO foo(name, bio) values ($1, $2) returning name
`

type InsertAtParamsParams struct {
	Name string
	Bio  string
}

func (q *Queries) InsertAtParams(ctx context.Context, arg InsertAtParamsParams) (string, error) {
	row := q.db.QueryRow(ctx, insertAtParams, arg.Name, arg.Bio)
	var name string
	err := row.Scan(&name)
	return name, err
}

const insertFuncParams = `-- name: InsertFuncParams :one
INSERT INTO foo(name, bio) values ($1, $2) returning name
`

type InsertFuncParamsParams struct {
	Name string
	Bio  string
}

func (q *Queries) InsertFuncParams(ctx context.Context, arg InsertFuncParamsParams) (string, error) {
	row := q.db.QueryRow(ctx, insertFuncParams, arg.Name, arg.Bio)
	var name string
	err := row.Scan(&name)
	return name, err
}

const update = `-- name: Update :one
UPDATE foo
SET
  name = CASE WHEN $1::bool
    THEN $2::text
    ELSE name
    END
RETURNING name, bio
`

type UpdateParams struct {
	SetName bool
	Name    string
}

func (q *Queries) Update(ctx context.Context, arg UpdateParams) (Foo, error) {
	row := q.db.QueryRow(ctx, update, arg.SetName, arg.Name)
	var i Foo
	err := row.Scan(&i.Name, &i.Bio)
	return i, err
}
