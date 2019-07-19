# Returning rows

To generate a database access method, annotate a query with a specific comment.

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL,
  birth_year int    NOT NULL
);


-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY id;
```

A few new pieces of code are generated beyond the `Author` struct. An interface
for the underlying database is generated. The `*sql.DB` and `*sql.Tx` types
satisty this interface.

The database access methods are added to a `Queries` struct, which is created
using the `New` method.

Note that the `*` in our query has been replaced with explicit column names.
This change ensures that the query will never return unexpected data.

Our query was annotated with `:one`, meaning that it should only return a
single row. We scan the data from that one into a `Author` struct.

Since the get query has a single parameter, the `GetAuthor` method takes a single
`int` as an argument.

Since the list query has no parameters, the `ListAuthors` method accepts no
arguments.


```go
package db

import (
	"context"
	"database/sql"
)

type Author struct {
	ID        int
	Bio       string
	BirthYear int
}

type dbtx interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db dbtx) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db dbtx
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, bio, birth_year FROM authors
WHERE id = $1
`

func (q *Queries) GetAuthor(ctx context.Context, id int) (Author, error) {
  row := q.db.QueryRowContext(ctx, getAuthor, id)
  var i Author
  err := row.Scan(&i.ID, &i.Bio, &i.BirthYear)
  return i, err
}

const listAuthors = `-- name: ListAuthors :many
SELECT id, bio, birth_year FROM authors
ORDER BY id
`

func (q *Queries) ListAuthors(ctx context.Context) ([]Author, error) {
  rows, err := q.db.QueryContext(ctx, listAuthors)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  var items []Author
  for rows.Next() {
    var i Author
    if err := rows.Scan(&i.ID, &i.Bio, &i.BirthYear); err != nil {
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
```

# Returning specific columns

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL,
  birth_year int    NOT NULL
);

-- name: GetBioForAuthor :one
SELECT bio FROM authors
WHERE id = $1;

-- name: GetInfoForAuthor :one
SELECT bio, birth_year FROM authors
WHERE id = $1;
```

When selecting a single column, only that value that returned. The `GetBioForAuthor`
method takes a single `int` as an argument and returns a `string` and an
`error`.

When selecting multiple columns, a row record (method-specific struct) is
returned. In this case, `GetInfoForAuthor` returns a struct with two fields:
`Bio` and `BirthYear`.

```go
package db

import (
	"context"
	"database/sql"
)

type dbtx interface {
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db dbtx) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db dbtx
}

const getBioForAuthor = `-- name: GetBioForAuthor :one
SELECT bio FROM authors
WHERE id = $1
`

func (q *Queries) GetBioForAuthor(ctx context.Context, id int) (string, error) {
	row := q.db.QueryRowContext(ctx, getBioForAuthor, id)
	var i string
	err := row.Scan(&i)
	return i, err
}

const getInfoForAuthor = `-- name: GetInfoForAuthor :one
SELECT bio, birth_year FROM authors
WHERE id = $1
`

type GetInfoForAuthorRow struct {
	Bio       string
	BirthYear int
}

func (q *Queries) GetBioForAuthor(ctx context.Context, id int) (GetBioForAuthor, error) {
	row := q.db.QueryRowContext(ctx, getInfoForAuthor, id)
	var i GetBioForAuthor
	err := row.Scan(&i.Bio, &i.BirthYear)
	return i, err
}
```
