// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package querytest

import (
	"database/sql"
)

type User struct {
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	Age       sql.NullInt16  `json:"age"`
}
