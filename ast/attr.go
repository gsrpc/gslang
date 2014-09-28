package ast

//Attr AST attribute node
type Attr struct {
	BasicExpr      //Mixin basic expr implement
	Type      Expr //attribute type ref
	Args      Expr //the argument list,maybe nil
}

//NewAttr create new attr node
func (node *Script) NewAttr(attrType Expr) *Attr {
	expr := &Attr{
		Type: attrType,
	}

	expr.Init(attrType.Name(), node)

	expr.Type.SetParent(expr)

	return expr
}
