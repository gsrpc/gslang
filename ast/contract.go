package ast

import (
	"fmt"

	"github.com/gsdocker/gserrors"
)

//Param AST method param node
type Param struct {
	BasicExpr      //Mixin basic expr implement
	ID        int  //param index
	Type      Expr // Param type
}

//Method AST method node
type Method struct {
	BasicExpr          //Mixin basic expr implement
	Return    []*Param //Return param
	Params    []*Param //Input parameters
}

//NewReturn create new return param
func (method *Method) NewReturn(paramType Expr) *Param {
	param := &Param{
		ID:   len(method.Return),
		Type: paramType,
	}
	paramType.SetParent(param)
	param.Init(fmt.Sprintf("return_arg(%d)", len(method.Return)), method.Script())
	param.SetParent(method)

	method.Return = append(method.Return, param)

	return param
}

//NewParam create new return param
func (method *Method) NewParam(paramType Expr) *Param {
	param := &Param{
		ID:   len(method.Params),
		Type: paramType,
	}
	paramType.SetParent(param)
	param.Init(fmt.Sprintf("arg(%d)", len(method.Params)), method.Script())
	param.SetParent(method)

	method.Params = append(method.Params, param)

	return param
}

//Contract AST contract node
type Contract struct {
	BasicExpr                    //Mixin basic expr implement
	Methods   map[string]*Method //The method table belong to contract
	Bases     []*TypeRef         //the contract's base contract list
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

//NewBase create new base table
func (expr *Contract) NewBase(base *TypeRef) (ref *TypeRef, ok bool) {
	for _, old := range expr.Bases {
		if base.Name() == old.Name() {
			ref = old
			return
		}
	}
	base.SetParent(expr)
	expr.Bases = append(expr.Bases, base)
	ref = base
	ok = true
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
