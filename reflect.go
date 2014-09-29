package gslang

import (
	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslang/ast"
)

//Enum get enum values
func Enum(enum *ast.Enum) map[string]int64 {
	values := make(map[string]int64)
	for _, val := range enum.Values {
		values[val.Name()] = val.Value
	}

	return values
}

//EvalFieldInitArg get table field init arg
func EvalFieldInitArg(field *ast.Field, expr ast.Expr) (ast.Expr, bool) {
	eval := &evalArg{
		field: field,
	}

	expr.Accept(eval)

	if eval.expr != nil {
		return eval.expr, true
	}

	return nil, false
}

//EvalEnumVal eval const enum expr
func EvalEnumVal(expr ast.Expr) int64 {
	visitor := &evalEnumVal{}
	expr.Accept(visitor)
	return visitor.val
}

//IsAttrUsage check if the expr is AttrUsage table
func IsAttrUsage(expr *ast.Table) bool {
	if expr.Package().Name() == GSLangPackage &&
		expr.Name() == "AttrUsage" {
		return true

	}

	return false
}

//EvalAttrUsage get attribute's usage attr val
func (cs *CompileS) EvalAttrUsage(attr *ast.Attr) int64 {

	gserrors.Require(attr.Type.Ref != nil, "attr(%s) must linked first :\n\t%s", attr, Pos(attr))

	table, ok := attr.Type.Ref.(*ast.Table)

	if !ok {
		gserrors.Panicf(
			ErrCompileS,
			"only table can be used as attribute type :\n\tattr def :%s\n\ttype def:",
			Pos(attr),
			Pos(attr.Type.Ref),
		)
	}

	for _, metattr := range table.Attrs() {
		usage, ok := metattr.Type.Ref.(*ast.Table)

		gserrors.Require(ok, "attr(%s) must linked first :\n\t%s", metattr, Pos(attr))

		if IsAttrUsage(usage) {
			field := usage.Fields["Target"]

			if field == nil {
				gserrors.Panicf(
					ErrCompileS,
					"inner error: gslang AttrUsage must declare Target Field \n\ttype def:",
					Pos(usage),
				)
			}

			if target, ok := EvalFieldInitArg(field, metattr.Args); ok {
				return EvalEnumVal(target)
			}

			gserrors.Panicf(
				ErrCompileS,
				"AttrUsage attribute initlist expect target val \n\tattr def:",
				Pos(metattr),
			)
		}
	}

	gserrors.Panicf(
		ErrCompileS,
		"target table can't be used as attribute type :\n\tattr def :%s\n\ttype def:%s",
		Pos(attr),
		Pos(attr.Type.Ref),
	)

	return 0
}
