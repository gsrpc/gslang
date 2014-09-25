package ast

import "github.com/gsdocker/gserrors"

//Method AST method node
type Method struct {
	BasicExpr        //Mixin basic expr implement
	Return    Expr   //Return param
	Params    []Expr //Input parameters
}

//Contract AST contract node
type Contract struct {
	BasicExpr                    //Mixin basic expr implement
	Methods   map[string]*Method //The method table belong to contract
}

//NewContract create new contract node
func (node *Script) NewContract(name string) (expr *Contract) {

	defer gserrors.Ensure(func() bool {
		return expr.Methods != nil
	}, "make sure alloc Contract's Methods field")

	expr = &Contract{
		Methods: make(map[string]*Method),
	}

	expr.Init(name, node)

	return
}

//NewMethod create new method node belong to current contract
func (expr *Contract) NewMethod(name string) (method *Method, ok bool) {
	gserrors.Require(expr.Methods != nil, "NewContract method must alloc Methods field")

	defer gserrors.Ensure(func() bool {
		return ok == (method == expr.Methods[name])
	}, "post condition check")

	defer gserrors.Ensure(func() bool {
		return expr.Methods != nil
	}, "make sure alloc Contract's Methods field")

	method = &Method{}

	method.Init(name, expr.Script())

	method.SetParent(expr)

	expr.Methods[name] = method

	ok = true

	return
}
