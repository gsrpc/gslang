package ast

import "github.com/gsdocker/gserrors"

//EnumVal the enum val node
type EnumVal struct {
	BasicExpr       //mixin default expr implement
	Value     int64 //value
}

//Enum the enum node
type Enum struct {
	BasicExpr                     //mixin default expr implement
	Values    map[string]*EnumVal //enum values table
}

//NewEnum create new enum node object
func (node *Script) NewEnum(name string) (expr *Enum) {

	defer gserrors.Ensure(func() bool {
		return expr.Values != nil
	}, "make sure alloc Enum's Values field")

	expr = &Enum{
		Values: make(map[string]*EnumVal),
	}

	expr.Init(name, node)

	return expr
}

//NewVal create new Enum val node
func (expr *Enum) NewVal(name string, val int64) (result *EnumVal, ok bool) {

	defer gserrors.Ensure(func() bool {
		return ok == (result == expr.Values[name])
	}, "post condition check")

	if result, ok = expr.Values[name]; ok {
		ok = !ok
		return
	}

	result = &EnumVal{
		Value: val,
	}

	result.Init(name, expr.Script())
	result.SetParent(expr)
	expr.Values[name] = result
	ok = true
	return
}
