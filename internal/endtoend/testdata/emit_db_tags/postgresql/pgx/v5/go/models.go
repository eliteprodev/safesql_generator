// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package querytest

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        int32       `db:"id"`
	FirstName string      `db:"first_name"`
	LastName  pgtype.Text `db:"last_name"`
	Age       int32       `db:"age"`
}
