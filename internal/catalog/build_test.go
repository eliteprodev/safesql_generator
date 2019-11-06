package catalog

import (
	"strconv"
	"testing"

	"github.com/kyleconroy/sqlc/internal/pg"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	query "github.com/lfittl/pg_query_go"
)

func buildCatalog(stmt string) (pg.Catalog, error) {
	c := pg.NewCatalog()
	tree, err := query.Parse(stmt)
	if err != nil {
		return c, err
	}
	for _, stmt := range tree.Statements {
		if err := Update(&c, stmt); err != nil {
			return c, err
		}
	}
	return c, nil
}

func TestUpdate(t *testing.T) {
	for i, tc := range []struct {
		stmt string
		c    pg.Catalog
	}{
		{
			"CREATE TYPE status AS ENUM ('open', 'closed');",
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Enums: map[string]pg.Enum{
							"status": pg.Enum{
								Name: "status",
								Vals: []string{"open", "closed"},
							},
						},
					},
				},
			},
		},
		{
			"CREATE TABLE venues ();",
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Tables: map[string]pg.Table{
							"venues": pg.Table{
								Name: "venues",
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ADD COLUMN bar text;
			ALTER TABLE foo DROP COLUMN bar;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Tables: map[string]pg.Table{
							"foo": pg.Table{
								Name: "foo",
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo DROP COLUMN IF EXISTS bar;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Tables: map[string]pg.Table{
							"foo": pg.Table{
								Name: "foo",
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo (bar text);
			ALTER TABLE foo ALTER bar SET NOT NULL;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Tables: map[string]pg.Table{
							"foo": pg.Table{
								Name:    "foo",
								Columns: []pg.Column{{Name: "bar", DataType: "text", NotNull: true, Table: pg.FQN{Schema: "public", Rel: "foo"}}},
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo (bar text[] not null);
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Tables: map[string]pg.Table{
							"foo": pg.Table{
								Name: "foo",
								Columns: []pg.Column{
									{Name: "bar", DataType: "text", IsArray: true, NotNull: true, Table: pg.FQN{Schema: "public", Rel: "foo"}},
								},
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo (bar text NOT NULL);
			ALTER TABLE foo ALTER bar DROP NOT NULL;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Tables: map[string]pg.Table{
							"foo": pg.Table{
								Name:    "foo",
								Columns: []pg.Column{{Name: "bar", DataType: "text", Table: pg.FQN{Schema: "public", Rel: "foo"}}},
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo (bar text NOT NULL);
			ALTER TABLE foo ALTER COLUMN bar DROP NOT NULL;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Tables: map[string]pg.Table{
							"foo": pg.Table{
								Name:    "foo",
								Columns: []pg.Column{{Name: "bar", DataType: "text", Table: pg.FQN{Schema: "public", Rel: "foo"}}},
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo (bar text);
			ALTER TABLE foo RENAME bar TO baz;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Tables: map[string]pg.Table{
							"foo": pg.Table{
								Name:    "foo",
								Columns: []pg.Column{{Name: "baz", DataType: "text", Table: pg.FQN{Schema: "public", Rel: "foo"}}},
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TABLE foo (bar text);
			ALTER TABLE foo ALTER bar SET DATA TYPE bool;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Tables: map[string]pg.Table{
							"foo": pg.Table{
								Name:    "foo",
								Columns: []pg.Column{{Name: "bar", DataType: "bool", Table: pg.FQN{Schema: "public", Rel: "foo"}}},
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE SCHEMA foo;
			CREATE SCHEMA bar;
			CREATE TABLE foo.baz ();
			ALTER TABLE foo.baz SET SCHEMA bar;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {},
					"foo":    {},
					"bar": {
						Tables: map[string]pg.Table{
							"baz": pg.Table{
								Name: "baz",
							},
						},
					},
				},
			},
		},
		{
			"CREATE TYPE status AS ENUM ('open', 'closed');",
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Enums: map[string]pg.Enum{
							"status": pg.Enum{
								Name: "status",
								Vals: []string{"open", "closed"},
							},
						},
						Tables: map[string]pg.Table{},
					},
				},
			},
		},
		{
			`
			CREATE TABLE venues ();
			DROP TABLE venues;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TABLE venues ();
			DROP TABLE IF EXISTS venues;
			DROP TABLE IF EXISTS venues;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TABLE venues ();
			ALTER TABLE venues RENAME TO arenas;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Enums: map[string]pg.Enum{},
						Tables: map[string]pg.Table{
							"arenas": pg.Table{
								Name: "arenas",
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE TYPE status AS ENUM ('open', 'closed');
			DROP TYPE status;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TYPE status AS ENUM ('open', 'closed');
			DROP TYPE IF EXISTS status;
			DROP TYPE IF EXISTS status;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TABLE venues ();
			DROP TABLE public.venues;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TYPE status AS ENUM ('open', 'closed');
			DROP TYPE public.status;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TYPE status AS ENUM ('open', 'closed');
			DROP TYPE public.status;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE SCHEMA foo;
			DROP SCHEMA foo;
			`,
			pg.NewCatalog(),
		},
		{
			`
			DROP SCHEMA IF EXISTS foo;
			`,
			pg.NewCatalog(),
		},
		{
			`
			CREATE TABLE venues (id SERIAL PRIMARY KEY);
			ALTER TABLE venues DROP CONSTRAINT venues_id_pkey;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Enums: map[string]pg.Enum{},
						Tables: map[string]pg.Table{
							"venues": pg.Table{
								Name: "venues",
								Columns: []pg.Column{
									{Name: "id", DataType: "serial", NotNull: true, Table: pg.FQN{Schema: "public", Rel: "venues"}},
								},
							},
						},
					},
				},
			},
		},
		{ // same name, different arity
			`
			CREATE FUNCTION foo(bar TEXT) RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
			CREATE FUNCTION foo(bar TEXT, baz TEXT) RETURNS TEXT AS $$ SELECT "baz" $$ LANGUAGE sql;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Funcs: map[string][]pg.Function{
							"foo": []pg.Function{
								{
									Name: "foo",
									Arguments: []pg.Argument{
										{
											Name:     "bar",
											DataType: "text",
										},
									},
									ReturnType: "bool",
								},
								{
									Name: "foo",
									Arguments: []pg.Argument{
										{
											Name:     "bar",
											DataType: "text",
										},
										{
											Name:     "baz",
											DataType: "text",
										},
									},
									ReturnType: "text",
								},
							},
						},
					},
				},
			},
		},
		{ // same name and arity, different arg types
			`
			CREATE FUNCTION foo(bar TEXT) RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
			CREATE FUNCTION foo(bar INTEGER) RETURNS TEXT AS $$ SELECT "baz" $$ LANGUAGE sql;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Funcs: map[string][]pg.Function{
							"foo": []pg.Function{
								{
									Name: "foo",
									Arguments: []pg.Argument{
										{
											Name:     "bar",
											DataType: "text",
										},
									},
									ReturnType: "bool",
								},
								{
									Name: "foo",
									Arguments: []pg.Argument{
										{
											Name:     "bar",
											DataType: "pg_catalog.int4",
										},
									},
									ReturnType: "text",
								},
							},
						},
					},
				},
			},
		},
		{
			`
			CREATE FUNCTION foo(bar TEXT, baz TEXT="baz") RETURNS bool AS $$ SELECT true $$ LANGUAGE sql;
			`,
			pg.Catalog{
				Schemas: map[string]pg.Schema{
					"public": {
						Funcs: map[string][]pg.Function{
							"foo": []pg.Function{
								{
									Name: "foo",
									Arguments: []pg.Argument{
										{
											Name:     "bar",
											DataType: "text",
										},
										{
											Name:       "baz",
											DataType:   "text",
											HasDefault: true,
										},
									},
									ReturnType: "bool",
								},
							},
						},
					},
				},
			},
		},
	} {
		test := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			c, err := buildCatalog(test.stmt)
			if err != nil {
				t.Log(test.stmt)
				t.Fatal(err)
			}

			delete(c.Schemas, "pg_catalog")
			delete(test.c.Schemas, "pg_catalog")

			if diff := cmp.Diff(test.c, c, cmpopts.EquateEmpty()); diff != "" {
				t.Log(test.stmt)
				t.Errorf("catalog mismatch:\n%s", diff)
			}
		})
	}
}

func TestUpdateErrors(t *testing.T) {
	for i, tc := range []struct {
		stmt string
		err  pg.Error
	}{
		{
			`
			CREATE TABLE foo ();
			CREATE TABLE foo ();
			`,
			pg.Error{Code: "42P07", Message: "relation \"foo\" already exists"},
		},
		{
			`
			CREATE TYPE foo AS ENUM ('bar');
			CREATE TYPE foo AS ENUM ('bar');
			`,
			pg.Error{Code: "42710", Message: "type \"foo\" already exists"},
		},
		{
			`
			DROP TABLE foo;
			`,
			pg.Error{Code: "42P01", Message: "relation \"foo\" does not exist"},
		},
		{
			`
			DROP TYPE foo;
			`,
			pg.Error{Code: "42704", Message: "type \"foo\" does not exist"},
		},
		{
			`
			CREATE TABLE foo ();
			CREATE TABLE bar ();
			ALTER TABLE foo RENAME TO bar;
			`,
			pg.Error{Code: "42P07", Message: "relation \"bar\" already exists"},
		},
		{
			`
			ALTER TABLE foo RENAME TO bar;
			`,
			pg.Error{Code: "42P01", Message: "relation \"foo\" does not exist"},
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ADD COLUMN bar text;
			ALTER TABLE foo ADD COLUMN bar text;
			`,
			pg.Error{Code: "42701", Message: "column \"bar\" of relation \"foo\" already exists"},
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo DROP COLUMN bar;
			`,
			pg.Error{Code: "42703", Message: "column \"bar\" of relation \"foo\" does not exist"},
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ALTER COLUMN bar SET NOT NULL;
			`,
			pg.Error{Code: "42703", Message: "column \"bar\" of relation \"foo\" does not exist"},
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ALTER COLUMN bar DROP NOT NULL;
			`,
			pg.Error{Code: "42703", Message: "column \"bar\" of relation \"foo\" does not exist"},
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo ALTER COLUMN bar DROP NOT NULL;
			`,
			pg.Error{Code: "42703", Message: "column \"bar\" of relation \"foo\" does not exist"},
		},
		{
			`
			CREATE SCHEMA foo;
			CREATE SCHEMA foo;
			`,
			pg.Error{Code: "42P06", Message: "schema \"foo\" already exists"},
		},
		{
			`
			ALTER TABLE foo.baz SET SCHEMA bar;
			`,
			pg.Error{Code: "3F000", Message: "schema \"foo\" does not exist"},
		},
		{
			`
			CREATE SCHEMA foo;
			ALTER TABLE foo.baz SET SCHEMA bar;
			`,
			pg.Error{Code: "42P01", Message: "relation \"baz\" does not exist"},
		},
		{
			`
			CREATE SCHEMA foo;
			CREATE TABLE foo.baz ();
			ALTER TABLE foo.baz SET SCHEMA bar;
			`,
			pg.Error{Code: "3F000", Message: "schema \"bar\" does not exist"},
		},
		{
			`
			DROP SCHEMA bar;
			`,
			pg.Error{Code: "3F000", Message: "schema \"bar\" does not exist"},
		},
		{
			`
			ALTER TABLE foo RENAME bar TO baz;
			`,
			pg.Error{Code: "42P01", Message: "relation \"foo\" does not exist"},
		},
		{
			`
			CREATE TABLE foo ();
			ALTER TABLE foo RENAME bar TO baz;
			`,
			pg.Error{Code: "42703", Message: "column \"bar\" of relation \"foo\" does not exist"},
		},
		{
			`
			CREATE TABLE foo (bar text, baz text);
			ALTER TABLE foo RENAME bar TO baz;
			`,
			pg.Error{Code: "42701", Message: "column \"baz\" of relation \"foo\" already exists"},
		},
	} {
		test := tc
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := buildCatalog(test.stmt)
			if err == nil {
				t.Log(test.stmt)
				t.Fatal("err was nil")
			}

			var actual pg.Error
			if err != nil {
				pge, ok := err.(pg.Error)
				if !ok {
					t.Log(test.stmt)
					t.Fatal(err)
				}
				actual = pge
			}

			if diff := cmp.Diff(test.err, actual, cmpopts.EquateEmpty()); diff != "" {
				t.Log(test.stmt)
				t.Errorf("error mismatch: \n%s", diff)
			}
		})
	}
}
