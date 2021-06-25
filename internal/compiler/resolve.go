package compiler

import (
	"fmt"
	"strconv"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

func dataType(n *ast.TypeName) string {
	if n.Schema != "" {
		return n.Schema + "." + n.Name
	} else {
		return n.Name
	}
}

func resolveCatalogRefs(c *catalog.Catalog, qc *QueryCatalog, rvs []*ast.RangeVar, args []paramRef, names map[int]string) ([]Parameter, error) {
	aliasMap := map[string]*ast.TableName{}
	// TODO: Deprecate defaultTable
	var defaultTable *ast.TableName
	var tables []*ast.TableName

	parameterName := func(n int, defaultName string) string {
		if n, ok := names[n]; ok {
			return n
		}
		return defaultName
	}

	typeMap := map[string]map[string]map[string]*catalog.Column{}
	indexTable := func(table catalog.Table) error {
		tables = append(tables, table.Rel)
		if defaultTable == nil {
			defaultTable = table.Rel
		}
		schema := table.Rel.Schema
		if schema == "" {
			schema = c.DefaultSchema
		}
		if _, exists := typeMap[schema]; !exists {
			typeMap[schema] = map[string]map[string]*catalog.Column{}
		}
		typeMap[schema][table.Rel.Name] = map[string]*catalog.Column{}
		for _, c := range table.Columns {
			cc := c
			typeMap[schema][table.Rel.Name][c.Name] = cc
		}
		return nil
	}

	for _, rv := range rvs {
		if rv.Relname == nil {
			continue
		}
		fqn, err := ParseTableName(rv)
		if err != nil {
			return nil, err
		}
		table, err := c.GetTable(fqn)
		if err != nil {
			// If the table name doesn't exist, fisrt check if it's a CTE
			if _, qcerr := qc.GetTable(fqn); qcerr != nil {
				return nil, err
			}
			continue
		}
		err = indexTable(table)
		if err != nil {
			return nil, err
		}
		if rv.Alias != nil {
			aliasMap[*rv.Alias.Aliasname] = fqn
		}
	}

	var a []Parameter
	for _, ref := range args {
		switch n := ref.parent.(type) {

		case *limitOffset:
			a = append(a, Parameter{
				Number: ref.ref.Number,
				Column: &Column{
					Name:     parameterName(ref.ref.Number, "offset"),
					DataType: "integer",
					NotNull:  true,
				},
			})

		case *limitCount:
			a = append(a, Parameter{
				Number: ref.ref.Number,
				Column: &Column{
					Name:     parameterName(ref.ref.Number, "limit"),
					DataType: "integer",
					NotNull:  true,
				},
			})

		case *ast.A_Expr:
			// TODO: While this works for a wide range of simple expressions,
			// more complicated expressions will cause this logic to fail.
			list := astutils.Search(n.Lexpr, func(node ast.Node) bool {
				_, ok := node.(*ast.ColumnRef)
				return ok
			})

			if len(list.Items) == 0 {
				// TODO: Move this to database-specific engine package
				dataType := "any"
				if astutils.Join(n.Name, ".") == "||" {
					dataType = "string"
				}
				a = append(a, Parameter{
					Number: ref.ref.Number,
					Column: &Column{
						Name:     parameterName(ref.ref.Number, ""),
						DataType: dataType,
					},
				})
				continue
			}

			switch left := list.Items[0].(type) {
			case *ast.ColumnRef:
				items := stringSlice(left.Fields)
				var key, alias string
				switch len(items) {
				case 1:
					key = items[0]
				case 2:
					alias = items[0]
					key = items[1]
				default:
					panic("too many field items: " + strconv.Itoa(len(items)))
				}

				search := tables
				if alias != "" {
					if original, ok := aliasMap[alias]; ok {
						search = []*ast.TableName{original}
					} else {
						for _, fqn := range tables {
							if fqn.Name == alias {
								search = []*ast.TableName{fqn}
							}
						}
					}
				}

				var found int
				for _, table := range search {
					schema := table.Schema
					if schema == "" {
						schema = c.DefaultSchema
					}
					if c, ok := typeMap[schema][table.Name][key]; ok {
						found += 1
						if ref.name != "" {
							key = ref.name
						}
						a = append(a, Parameter{
							Number: ref.ref.Number,
							Column: &Column{
								Name:     parameterName(ref.ref.Number, key),
								DataType: dataType(&c.Type),
								NotNull:  c.IsNotNull,
								IsArray:  c.IsArray,
								Length:   c.Length,
								Table:    table,
							},
						})
					}
				}

				if found == 0 {
					return nil, &sqlerr.Error{
						Code:     "42703",
						Message:  fmt.Sprintf("column \"%s\" does not exist", key),
						Location: left.Location,
					}
				}
				if found > 1 {
					return nil, &sqlerr.Error{
						Code:     "42703",
						Message:  fmt.Sprintf("column reference \"%s\" is ambiguous", key),
						Location: left.Location,
					}
				}
			}

		case *ast.FuncCall:
			fun, err := c.ResolveFuncCall(n)
			if err != nil {
				// Synthesize a function on the fly to avoid returning with an error
				// for an unknown Postgres function (e.g. defined in an extension)
				var args []*catalog.Argument
				for range n.Args.Items {
					args = append(args, &catalog.Argument{
						Type: &ast.TypeName{Name: "any"},
					})
				}
				fun = &catalog.Function{
					Name:       n.Func.Name,
					Args:       args,
					ReturnType: &ast.TypeName{Name: "any"},
				}
			}
			for i, item := range n.Args.Items {
				funcName := fun.Name
				var argName string
				switch inode := item.(type) {
				case *ast.ParamRef:
					if inode.Number != ref.ref.Number {
						continue
					}
				case *ast.TypeCast:
					pr, ok := inode.Arg.(*ast.ParamRef)
					if !ok {
						continue
					}
					if pr.Number != ref.ref.Number {
						continue
					}
				case *ast.NamedArgExpr:
					pr, ok := inode.Arg.(*ast.ParamRef)
					if !ok {
						continue
					}
					if pr.Number != ref.ref.Number {
						continue
					}
					if inode.Name != nil {
						argName = *inode.Name
					}
				default:
					continue
				}

				if fun.Args == nil {
					defaultName := funcName
					if argName != "" {
						defaultName = argName
					}
					a = append(a, Parameter{
						Number: ref.ref.Number,
						Column: &Column{
							Name:     parameterName(ref.ref.Number, defaultName),
							DataType: "any",
						},
					})
					continue
				}

				var paramName string
				var paramType *ast.TypeName
				if argName == "" {
					paramName = fun.Args[i].Name
					paramType = fun.Args[i].Type
				} else {
					paramName = argName
					for _, arg := range fun.Args {
						if arg.Name == argName {
							paramType = arg.Type
						}
					}
					if paramType == nil {
						panic(fmt.Sprintf("named argument %s has no type", paramName))
					}
				}
				if paramName == "" {
					paramName = funcName
				}

				a = append(a, Parameter{
					Number: ref.ref.Number,
					Column: &Column{
						Name:     parameterName(ref.ref.Number, paramName),
						DataType: dataType(paramType),
						NotNull:  true,
					},
				})
			}

			if fun.ReturnType == nil {
				continue
			}

			table, err := c.GetTable(&ast.TableName{
				Catalog: fun.ReturnType.Catalog,
				Schema:  fun.ReturnType.Schema,
				Name:    fun.ReturnType.Name,
			})
			if err != nil {
				// The return type wasn't a table.
				continue
			}
			err = indexTable(table)
			if err != nil {
				return nil, err
			}

		case *ast.ResTarget:
			if n.Name == nil {
				return nil, fmt.Errorf("*ast.ResTarget has nil name")
			}
			key := *n.Name

			var schema, rel string
			// TODO: Deprecate defaultTable
			if defaultTable != nil {
				schema = defaultTable.Schema
				rel = defaultTable.Name
			}
			if ref.rv != nil {
				fqn, err := ParseTableName(ref.rv)
				if err != nil {
					return nil, err
				}
				schema = fqn.Schema
				rel = fqn.Name
			}
			if schema == "" {
				schema = c.DefaultSchema
			}

			tableMap, ok := typeMap[schema][rel]
			if !ok {
				return nil, sqlerr.RelationNotFound(rel)
			}

			if c, ok := tableMap[key]; ok {
				a = append(a, Parameter{
					Number: ref.ref.Number,
					Column: &Column{
						Name:     parameterName(ref.ref.Number, key),
						DataType: dataType(&c.Type),
						NotNull:  c.IsNotNull,
						IsArray:  c.IsArray,
						Table:    &ast.TableName{Schema: schema, Name: rel},
						Length:   c.Length,
					},
				})
			} else {
				return nil, &sqlerr.Error{
					Code:     "42703",
					Message:  fmt.Sprintf("column \"%s\" does not exist", key),
					Location: n.Location,
				}
			}

		case *ast.TypeCast:
			if n.TypeName == nil {
				return nil, fmt.Errorf("*ast.TypeCast has nil type name")
			}
			col := toColumn(n.TypeName)
			col.Name = parameterName(ref.ref.Number, col.Name)
			a = append(a, Parameter{
				Number: ref.ref.Number,
				Column: col,
			})

		case *ast.ParamRef:
			a = append(a, Parameter{Number: ref.ref.Number})

		default:
			fmt.Printf("unsupported reference type: %T", n)
		}
	}
	return a, nil
}
