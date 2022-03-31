// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: query.sql

package querytest

import (
	"context"
	"database/sql"
)

const fooByBarB = `-- name: FooByBarB :many
SELECT a, b from foo where foo.a in (select a from bar where bar.b = ?)
`

func (q *Queries) FooByBarB(ctx context.Context, b sql.NullString) ([]Foo, error) {
	rows, err := q.db.QueryContext(ctx, fooByBarB, b)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.A, &i.B); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fooByList = `-- name: FooByList :many
SELECT a, b from foo where foo.a in (?, ?)
`

type FooByListParams struct {
	A   sql.NullString
	A_2 sql.NullString
}

func (q *Queries) FooByList(ctx context.Context, arg FooByListParams) ([]Foo, error) {
	rows, err := q.db.QueryContext(ctx, fooByList, arg.A, arg.A_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.A, &i.B); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fooByNotList = `-- name: FooByNotList :many
SELECT a, b from foo where foo.a not in (?, ?)
`

type FooByNotListParams struct {
	A   sql.NullString
	A_2 sql.NullString
}

func (q *Queries) FooByNotList(ctx context.Context, arg FooByNotListParams) ([]Foo, error) {
	rows, err := q.db.QueryContext(ctx, fooByNotList, arg.A, arg.A_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.A, &i.B); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fooByParamList = `-- name: FooByParamList :many
SELECT a, b from foo where ? in (foo.a, foo.b)
`

func (q *Queries) FooByParamList(ctx context.Context, a sql.NullString) ([]Foo, error) {
	rows, err := q.db.QueryContext(ctx, fooByParamList, a)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Foo
	for rows.Next() {
		var i Foo
		if err := rows.Scan(&i.A, &i.B); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
