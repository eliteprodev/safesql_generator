package ast

import ()

type CreateRangeStmt struct {
	TypeName *List
	Params   *List
}

func (n *CreateRangeStmt) Pos() int {
	return 0
}
