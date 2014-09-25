package ast

//Attr AST attribute node
type Attr struct {
	BasicExpr          //Mixin basic expr implement
	Type      *TypeRef //attribute type ref
}

//NewAttr create new attr node
func (node *Script) NewAttr(name string) *Attr {
	expr := &Attr{
		Type: node.NewTypeRef(name),
	}

	expr.Init(name, node)

	expr.Type.SetParent(expr)

	return expr
}
