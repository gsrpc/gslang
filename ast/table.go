package ast

//Field AST field node type
type Field struct {
	BasicExpr        //Inher default node implement
	ID        uint16 //Schema define field index
	Type      Expr   //field type
}

//Table AST table type node
type Table struct {
	BasicExpr          //Inher default expr implement
	Fields    []*Field //field table
}

//NewTable create new table node
func (script *Script) NewTable(name string) (expr *Table) {

	expr = &Table{}

	expr.Init(name, script)

	return
}

// Field get field
func (expr *Table) Field(name string) (*Field, bool) {
	for _, field := range expr.Fields {
		if field.Name() == name {
			return field, true
		}
	}

	return nil, false
}

//NewField create new
func (expr *Table) NewField(name string) (*Field, bool) {

	for _, field := range expr.Fields {
		if field.Name() == name {
			return field, false
		}
	}

	field := &Field{
		ID: uint16(len(expr.Fields)),
	}

	field.Init(name, expr.Script())

	field.SetParent(expr)

	expr.Fields = append(expr.Fields, field)

	return field, true
}
