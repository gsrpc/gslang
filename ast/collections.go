package ast

//Array AST array node
type Array struct {
	BasicExpr          //Mixin basic expr implement
	Component *TypeRef //array type ref
}

//NewArray create new array node
func (node *Script) NewArray(name string) *Array {
	expr := &Array{
		Component: node.NewTypeRef(name),
	}

	expr.Init(name, node)

	expr.Component.SetParent(expr)

	return expr
}
