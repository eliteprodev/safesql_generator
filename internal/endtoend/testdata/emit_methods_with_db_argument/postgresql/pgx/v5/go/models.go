// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package querytest

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID        int32
	FirstName string
	LastName  pgtype.Text
	Age       int32
}
