// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: query.sql

package querytest

import (
	"context"
)

const columnAs = `-- name: ColumnAs :many
SELECT email AS id FROM foo
`

func (q *Queries) ColumnAs(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, columnAs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
