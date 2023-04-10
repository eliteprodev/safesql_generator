package ast

type DeleteStmt struct {
	Relations     *List
	UsingClause   *List
	WhereClause   Node
	ReturningList *List
	WithClause    *WithClause
}

func (n *DeleteStmt) Pos() int {
	return 0
}
