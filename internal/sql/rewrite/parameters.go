package rewrite

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/source"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
	"github.com/kyleconroy/sqlc/internal/sql/named"
)

// Given an AST node, return the string representation of names
func flatten(root ast.Node) (string, bool) {
	sw := &stringWalker{}
	astutils.Walk(sw, root)
	return sw.String, sw.IsConst
}

type stringWalker struct {
	String  string
	IsConst bool
}

func (s *stringWalker) Visit(node ast.Node) astutils.Visitor {
	if _, ok := node.(*ast.A_Const); ok {
		s.IsConst = true
	}
	if n, ok := node.(*ast.String); ok {
		s.String += n.Str
	}
	return s
}

func isNamedParamSignCast(node ast.Node) bool {
	expr, ok := node.(*ast.A_Expr)
	if !ok {
		return false
	}
	_, cast := expr.Rexpr.(*ast.TypeCast)
	return astutils.Join(expr.Name, ".") == "@" && cast
}

// paramFromFuncCall creates a param from sqlc.n?arg() calls return the
// parameter and whether the parameter name was specified a best guess as its
// "source" string representation (used for replacing this function call in the
// original SQL query)
func paramFromFuncCall(call *ast.FuncCall) (named.Param, string) {
	paramName, isConst := flatten(call.Args)

	// origName keeps track of how the parameter was specified in the source SQL
	origName := paramName
	if isConst {
		origName = fmt.Sprintf("'%s'", paramName)
	}

	param := named.NewParam(paramName)
	if call.Func.Name == "narg" {
		param = named.NewUserNullableParam(paramName)
	}

	// TODO: This code assumes that sqlc.arg(name) / sqlc.narg(name) is on a single line
	// with no extraneous spaces (or any non-significant tokens for that matter)
	origText := fmt.Sprintf("%s.%s(%s)", call.Func.Schema, call.Func.Name, origName)
	return param, origText
}

func NamedParameters(engine config.Engine, raw *ast.RawStmt, numbs map[int]bool, dollar bool) (*ast.RawStmt, *named.ParamSet, []source.Edit) {
	foundFunc := astutils.Search(raw, named.IsParamFunc)
	foundSign := astutils.Search(raw, named.IsParamSign)
	hasNamedParameterSupport := engine != config.EngineMySQL
	allParams := named.NewParamSet(numbs, hasNamedParameterSupport)

	if len(foundFunc.Items)+len(foundSign.Items) == 0 {
		return raw, allParams, nil
	}

	var edits []source.Edit
	node := astutils.Apply(raw, func(cr *astutils.Cursor) bool {
		node := cr.Node()
		switch {
		case named.IsParamFunc(node):
			fun := node.(*ast.FuncCall)
			param, origText := paramFromFuncCall(fun)
			argn := allParams.Add(param)
			cr.Replace(&ast.ParamRef{
				Number:   argn,
				Location: fun.Location,
			})

			var replace string
			if engine == config.EngineMySQL || !dollar {
				replace = "?"
			} else {
				replace = fmt.Sprintf("$%d", argn)
			}

			edits = append(edits, source.Edit{
				Location: fun.Location - raw.StmtLocation,
				Old:      origText,
				New:      replace,
			})
			return false

		case isNamedParamSignCast(node):
			expr := node.(*ast.A_Expr)
			cast := expr.Rexpr.(*ast.TypeCast)
			paramName, _ := flatten(cast.Arg)
			param := named.NewParam(paramName)

			argn := allParams.Add(param)
			cast.Arg = &ast.ParamRef{
				Number:   argn,
				Location: expr.Location,
			}
			cr.Replace(cast)

			// TODO: This code assumes that @foo::bool is on a single line
			var replace string
			if engine == config.EngineMySQL || !dollar {
				replace = "?"
			} else {
				replace = fmt.Sprintf("$%d", argn)
			}

			edits = append(edits, source.Edit{
				Location: expr.Location - raw.StmtLocation,
				Old:      fmt.Sprintf("@%s", paramName),
				New:      replace,
			})
			return false

		case named.IsParamSign(node):
			expr := node.(*ast.A_Expr)
			paramName, _ := flatten(expr.Rexpr)
			param := named.NewParam(paramName)

			argn := allParams.Add(param)
			cr.Replace(&ast.ParamRef{
				Number:   argn,
				Location: expr.Location,
			})

			// TODO: This code assumes that @foo is on a single line
			var replace string
			if engine == config.EngineMySQL || !dollar {
				replace = "?"
			} else {
				replace = fmt.Sprintf("$%d", argn)
			}

			edits = append(edits, source.Edit{
				Location: expr.Location - raw.StmtLocation,
				Old:      fmt.Sprintf("@%s", paramName),
				New:      replace,
			})
			return false

		default:
			return true
		}
	}, nil)

	return node.(*ast.RawStmt), allParams, edits
}
