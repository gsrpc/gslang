package ast

//Array AST array node
type Array struct {
	BasicExpr        //Mixin basic expr implement
	Length    uint16 //array length
	Component Expr   //array type ref
}

//NewArray create new array node
func (node *Script) NewArray(length uint16, component Expr) *Array {
	expr := &Array{
		Length:    length,
		Component: component,
	}

	expr.Init(component.Name(), node)

	expr.Component.SetParent(expr)

	return expr
}

//List AST array node
type List struct {
	BasicExpr      //Mixin basic expr implement
	Component Expr //array type ref
}

//NewList create new array node
func (node *Script) NewList(component Expr) *List {
	expr := &List{
		Component: component,
	}

	expr.Init(component.Name(), node)

	expr.Component.SetParent(expr)

	return expr
}
