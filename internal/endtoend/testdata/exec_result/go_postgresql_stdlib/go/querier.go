// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.1

package querytest

import (
	"context"
	"database/sql"
)

type Querier interface {
	DeleteBarByID(ctx context.Context, id int32) (sql.Result, error)
}

var _ Querier = (*Queries)(nil)
