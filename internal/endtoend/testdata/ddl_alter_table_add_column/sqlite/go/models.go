// Code generated by sqlc. DO NOT EDIT.

package querytest

import (
	"database/sql"
)

type Venue struct {
	Name     sql.NullString
	Location sql.NullString
	Size     sql.NullInt64
}
