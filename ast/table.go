package ast

import "github.com/gsdocker/gserrors"

//Field AST field node type
type Field struct {
	BasicExpr        //Inher default node implement
	ID        uint16 //Schema define field index
	Type      Expr   //field type
}

//Table AST table type node
type Table struct {
	BasicExpr                   //Inher default expr implement
	Fields    map[string]*Field //field table
}

//NewTable create new table node
func (script *Script) NewTable(name string) (expr *Table) {

	defer gserrors.Ensure(func() bool {
		return expr.Fields != nil
	}, "make sure alloc Table's Fields field")

	expr = &Table{
		Fields: make(map[string]*Field),
	}

	expr.Init(name, script)

	return
}

//NewField create new
func (expr *Table) NewField(name string) (*Field, bool) {
	if old, ok := expr.Fields[name]; ok {
		return old, false
	}

	field := &Field{
		ID: uint16(len(expr.Fields)),
	}

	field.Init(name, expr.Script())

	field.SetParent(expr)

	expr.Fields[name] = field

	return field, true
}
