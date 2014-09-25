package ast

//TypeRef type reference node
type TypeRef struct {
	BasicExpr      //Mixin basic expr implement
	Ref       Expr //reference node,maybe : contract,enum table
}

//NewTypeRef create new type reference node
func (node *Script) NewTypeRef(name string) *TypeRef {
	expr := &TypeRef{}

	expr.Init(name, node)

	return expr
}
