// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.1
// source: queries.sql

package primary_key_later

import (
	"context"
)

const getAuthor = `-- name: GetAuthor :one
SELECT
    id, name, bio
FROM
    authors
WHERE
    id = $1
LIMIT 1
`

func (q *Queries) GetAuthor(ctx context.Context, id int64) (Author, error) {
	row := q.db.QueryRow(ctx, getAuthor, id)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}
