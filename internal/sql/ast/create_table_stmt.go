package ast

type CreateTableStmt struct {
	IfNotExists bool
	Name        *TableName
	Cols        []*ColumnDef
	ReferTable *TableName
}

func (n *CreateTableStmt) Pos() int {
	return 0
}
