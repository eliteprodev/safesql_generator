// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package querytest

import (
	"context"
)

const star = `-- name: Star :many
SELECT bid, fid FROM bar, foo
`

type StarRow struct {
	Bid int32
	Fid int32
}

func (q *Queries) Star(ctx context.Context) ([]StarRow, error) {
	rows, err := q.db.QueryContext(ctx, star)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []StarRow
	for rows.Next() {
		var i StarRow
		if err := rows.Scan(&i.Bid, &i.Fid); err != nil {
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
