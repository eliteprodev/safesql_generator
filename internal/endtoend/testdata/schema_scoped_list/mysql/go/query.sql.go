// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
)

const schemaScopedList = `-- name: SchemaScopedList :many
SELECT id FROM foo.bar
`

func (q *Queries) SchemaScopedList(ctx context.Context) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, schemaScopedList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
