package gslang

import (
	"bytes"
	"fmt"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslang/ast"
)

func (cs *CompileS) link(pkg *ast.Package) {
	linker := &Linker{
		CompileS: cs,
	}

	pkg.Accept(linker)

	linker2 := &attrLinker{
		CompileS: cs,
	}

	pkg.Accept(linker2)

	linker3 := &contractLinker{
		CompileS: cs,
	}

	pkg.Accept(linker3)
}

//Linker the type reference linker
type Linker struct {
	*CompileS        // compile services which this linker belongs to
	ast.EmptyVisitor //Mixin empty visitor implements
}

// VisitPackage implement visitor interface
func (linker *Linker) VisitPackage(pkg *ast.Package) ast.Node {

	for _, attr := range pkg.Attrs() {
		attr.Accept(linker)
	}

	for _, script := range pkg.Scripts {
		script.Accept(linker)
	}

	return pkg
}

// VisitScript implement visitor interface
func (linker *Linker) VisitScript(script *ast.Script) ast.Node {
	for _, attr := range script.Attrs() {
		attr.Accept(linker)
	}

	for _, expr := range script.Types {
		expr.Accept(linker)
	}

	return script
}

//VisitTable implement visitor interface
func (linker *Linker) VisitTable(table *ast.Table) ast.Node {
	for _, attr := range table.Attrs() {
		attr.Accept(linker)
	}

	for _, field := range table.Fields {
		field.Accept(linker)
	}

	return table
}

//VisitField implement visitor interface
func (linker *Linker) VisitField(field *ast.Field) ast.Node {
	for _, attr := range field.Attrs() {
		attr.Accept(linker)
	}

	field.Type.Accept(linker)

	return field
}

//VisitEnum implement visitor interface
func (linker *Linker) VisitEnum(enum *ast.Enum) ast.Node {
	for _, attr := range enum.Attrs() {
		attr.Accept(linker)
	}

	for _, val := range enum.Values {
		val.Accept(linker)
	}

	return enum
}

//VisitEnumVal implement visitor interface
func (linker *Linker) VisitEnumVal(val *ast.EnumVal) ast.Node {
	for _, attr := range val.Attrs() {
		attr.Accept(linker)
	}

	return val
}

//VisitContract implement visitor interface
func (linker *Linker) VisitContract(contract *ast.Contract) ast.Node {
	for _, attr := range contract.Attrs() {
		attr.Accept(linker)
	}

	for _, base := range contract.Bases {
		base.Accept(linker)
	}

	for _, method := range contract.Methods {
		method.Accept(linker)
	}

	return contract
}

//VisitMethod implement visitor interface
func (linker *Linker) VisitMethod(method *ast.Method) ast.Node {
	for _, attr := range method.Attrs() {
		attr.Accept(linker)
	}

	for _, expr := range method.Return {
		expr.Accept(linker)
	}

	for _, expr := range method.Params {
		expr.Accept(linker)
	}

	return method
}

//VisitParam implement visitor interface
func (linker *Linker) VisitParam(param *ast.Param) ast.Node {
	for _, attr := range param.Attrs() {
		attr.Accept(linker)
	}

	param.Type.Accept(linker)

	return param
}

//VisitBinaryOp implement visitor interface
func (linker *Linker) VisitBinaryOp(op *ast.BinaryOp) ast.Node {
	op.Left.Accept(linker)
	op.Right.Accept(linker)
	return op
}

//VisitList implement visitor interface
func (linker *Linker) VisitList(list *ast.List) ast.Node {
	list.Component.Accept(linker)
	return list
}

//VisitArray implement visitor interface
func (linker *Linker) VisitArray(array *ast.Array) ast.Node {
	array.Component.Accept(linker)
	return array
}

// VisitAttr implement visitor interface
func (linker *Linker) VisitAttr(attr *ast.Attr) ast.Node {
	attr.Type.Accept(linker)
	if attr.Args != nil {
		attr.Args.Accept(linker)
	}
	return attr
}

// VisitArgs implement visitor interface
func (linker *Linker) VisitArgs(args *ast.Args) ast.Node {
	for _, arg := range args.Items {
		arg.Accept(linker)
	}

	return args
}

// VisitNamedArgs implement visitor interface
func (linker *Linker) VisitNamedArgs(args *ast.NamedArgs) ast.Node {
	for _, arg := range args.Items {
		arg.Accept(linker)
	}

	return args
}

//VisitTypeRef implement visitor interface
func (linker *Linker) VisitTypeRef(ref *ast.TypeRef) ast.Node {
	if ref.Ref == nil {
		nodes := len(ref.NamePath)
		gserrors.Assert(nodes > 0, "the NamePath, can't be nil")

		switch nodes {
		case 1:
			if pkg, ok := ref.Script().Imports[ref.NamePath[0]]; !ok {
				pkg := ref.Package()
				gserrors.Assert(pkg != nil, "ref(%s) must bind ast tree", ref)
				if expr, ok := pkg.Types[ref.NamePath[0]]; ok {
					ref.Ref = expr
					return ref
				}
			} else {
				linker.errorf(
					Pos(ref),
					"type name conflict with import package name :\n\tsee:%s",
					Pos(pkg),
				)
			}
		case 2:
			if pkg, ok := ref.Script().Imports[ref.NamePath[0]]; ok {
				gserrors.Assert(
					pkg.Ref != nil,
					"(%s)first parse phase must link import pacakge:%s",
					ref.Script(),
					pkg,
				)
				if expr, ok := pkg.Ref.Types[ref.NamePath[1]]; ok {
					ref.Ref = expr
					return ref
				}
			} else {
				if expr, ok := ref.Package().Types[ref.NamePath[0]]; ok {
					if enum, ok := expr.(*ast.Enum); ok {
						if val, ok := enum.Values[ref.NamePath[1]]; ok {
							ref.Ref = val
							return ref
						}
					}
				}
			}
		case 3:
			if pkg, ok := ref.Script().Imports[ref.NamePath[0]]; ok {
				if expr, ok := pkg.Ref.Types[ref.NamePath[1]]; ok {
					if enum, ok := expr.(*ast.Enum); ok {
						if val, ok := enum.Values[ref.NamePath[2]]; ok {
							ref.Ref = val
							return ref
						}
					}
				}
			}
		}
	}
	linker.errorf(Pos(ref), "unknown type(%s)", ref)
	return ref
}

type contractLinker struct {
	*CompileS        // compile services which this linker belongs to
	ast.EmptyVisitor //Mixin empty visitor implements
}

// VisitPackage implement visitor interface
func (linker *contractLinker) VisitPackage(pkg *ast.Package) ast.Node {

	for _, script := range pkg.Scripts {
		script.Accept(linker)
	}

	return pkg
}

// VisitScript implement visitor interface
func (linker *contractLinker) VisitScript(script *ast.Script) ast.Node {
	for _, expr := range script.Types {
		expr.Accept(linker)
	}

	return script
}

//VisitContract implement visitor interface
func (linker *contractLinker) VisitContract(contract *ast.Contract) ast.Node {

	linker.unwind(contract, nil)

	return contract
}

func (linker *contractLinker) unwind(expr *ast.Contract, stack []*ast.Contract) []*ast.Contract {

	if _, ok := expr.Extra("unwind"); ok {
		return stack
	}

	var stream bytes.Buffer

	for _, contract := range stack {
		if contract == expr || stream.Len() != 0 {
			stream.WriteString(fmt.Sprintf("\t%s inheri\n", contract))
		}
	}

	if stream.Len() != 0 {
		linker.errorf(
			Pos(expr),
			"circular inheri :\n%s\t%s",
			stream.String(),
			expr,
		)
	}

	stack = append(stack, expr)

	modify := uint16(0)

	for _, base := range expr.Bases {
		contract, ok := base.Ref.(*ast.Contract)
		if !ok {
			linker.errorf(
				Pos(base),
				"contract(%s) inheri type is not contract :\n\tsee:%s",
				expr,
				Pos(base.Ref),
			)
		}

		stack = linker.unwind(contract, stack)

		modify = modify + uint16(len(contract.Methods))
	}

	for _, method := range expr.Methods {
		method.ID = method.ID + modify
	}

	modify = uint16(0)

	for _, base := range expr.Bases {
		contract := base.Ref.(*ast.Contract)

		for _, method := range contract.Methods {
			clone := &ast.Method{}
			*clone = *method
			clone.ID = clone.ID + modify

			if old, ok := expr.Methods[clone.Name()]; ok {
				linker.errorf(
					Pos(expr),
					"duplicate method name :%s\n\tsee:%s\n\tsee:%s",
					clone,
					Pos(old),
					Pos(clone),
				)
			}

			method.SetParent(expr)

			expr.Methods[clone.Name()] = clone
		}

		modify = modify + uint16(len(contract.Methods))
	}

	expr.NewExtra("unwind", true)

	stack = stack[:len(stack)-1]

	return stack
}

type attrLinker struct {
	*CompileS                         // compile services which this linker belongs to
	ast.EmptyVisitor                  //Mixin empty visitor implements
	attrTarget       map[string]int64 //attribute target enum map
	attrStruct       ast.Expr         //the struct attribute expr node
}

// VisitPackage implement visitor interface
func (linker *attrLinker) VisitPackage(pkg *ast.Package) ast.Node {

	if len(pkg.Scripts) == 0 {
		return pkg
	}

	if pkg.Name() == GSLangPackage {
		if expr, ok := pkg.Types[GSLangAttrTarget]; ok {
			if enum, ok := expr.(*ast.Enum); ok {
				linker.attrTarget = Enum(enum)
			}
		}
	} else {
		if pkg, ok := linker.Loaded[GSLangPackage]; ok {
			if expr, ok := pkg.Types[GSLangAttrTarget]; ok {
				if enum, ok := expr.(*ast.Enum); ok {
					linker.attrTarget = Enum(enum)
				}
			}
		}
	}

	if linker.attrTarget == nil {
		gserrors.Panicf(ErrCompileS, "inner error: can't found gslang.AttrTarget enum")
	}

	if pkg.Name() == GSLangPackage {

		linker.attrStruct = pkg.Types[GSLangAttrStruct]

		if linker.attrStruct == nil {
			gserrors.Panicf(ErrCompileS, "inner error: can't found gslang.Struct attribute type")
		}

	} else {
		attrStruct, err := linker.Type(GSLangPackage, GSLangAttrStruct)

		if err != nil {
			gserrors.Panicf(err, "inner error: can't found gslang.Struct attribute type")
		}

		linker.attrStruct = attrStruct
	}

	for _, script := range pkg.Scripts {
		script.Accept(linker)
	}

	return pkg
}

// VisitScript implement visitor interface
func (linker *attrLinker) VisitScript(script *ast.Script) ast.Node {
	for _, attr := range script.Attrs() {
		target := linker.EvalAttrUsage(attr)

		if target&linker.attrTarget["Script"] == 0 {

			if target&linker.attrTarget["Package"] != 0 {
				script.RemoveAttr(attr)
				script.Package().AddAttr(attr)

			} else {
				linker.errorf(
					Pos(attr),
					"attr(%s) can't be used to attribute script :\n\tsee:%s",
					attr,
					Pos(attr.Type.Ref),
				)
			}
		}
	}

	for _, expr := range script.Types {
		expr.Accept(linker)
	}

	return script
}

//VisitTable implement visitor interface
func (linker *attrLinker) VisitTable(table *ast.Table) ast.Node {

	var isStruct bool
	//detect if this table is an struct
	if len(ast.GetAttrs(table, linker.attrStruct)) > 0 {
		isStruct = true
		markAsStruct(table)
	}

	for _, attr := range table.Attrs() {
		target := linker.EvalAttrUsage(attr)

		var toMove bool

		if isStruct {

			if target&linker.attrTarget["Struct"] == 0 {
				toMove = true
			}
		} else {
			if target&linker.attrTarget["Table"] == 0 {
				toMove = true
			}
		}

		if toMove {
			if target&linker.attrTarget["Script"] != 0 {
				table.RemoveAttr(attr)
				table.Script().AddAttr(attr)
				continue
			}

			if target&linker.attrTarget["Package"] != 0 {
				table.RemoveAttr(attr)
				table.Package().AddAttr(attr)
				continue
			}
			linker.errorf(
				Pos(attr),
				"attr(%s) can't be used to attribute table/struct :\n\tsee:%s",
				attr,
				Pos(attr.Type.Ref),
			)
		}

	}

	for _, field := range table.Fields {
		field.Accept(linker)
	}

	return table
}

//VisitField implement visitor interface
func (linker *attrLinker) VisitField(field *ast.Field) ast.Node {
	for _, attr := range field.Attrs() {

		target := linker.EvalAttrUsage(attr)

		if target&linker.attrTarget["Field"] == 0 {
			linker.errorf(
				Pos(attr),
				"attr(%s) can't be used to attribute field :\n\tsee:%s",
				attr,
				Pos(attr.Type.Ref),
			)
		}
	}

	return field
}

//VisitEnum implement visitor interface
func (linker *attrLinker) VisitEnum(enum *ast.Enum) ast.Node {
	for _, attr := range enum.Attrs() {
		target := linker.EvalAttrUsage(attr)

		if target&linker.attrTarget["Enum"] == 0 {
			linker.errorf(
				Pos(attr),
				"attr(%s) can't be used to attribute enum :\n\tsee:%s",
				attr,
				Pos(attr.Type.Ref),
			)
		}
	}

	for _, val := range enum.Values {
		val.Accept(linker)
	}

	return enum
}

//VisitEnumVal implement visitor interface
func (linker *attrLinker) VisitEnumVal(val *ast.EnumVal) ast.Node {
	for _, attr := range val.Attrs() {
		target := linker.EvalAttrUsage(attr)

		if target&linker.attrTarget["EnumVal"] == 0 {
			linker.errorf(
				Pos(attr),
				"attr(%s) can't be used to attribute enum value :\n\tsee:%s",
				attr,
				Pos(attr.Type.Ref),
			)
		}
	}

	return val
}

//VisitContract implement visitor interface
func (linker *attrLinker) VisitContract(contract *ast.Contract) ast.Node {
	for _, attr := range contract.Attrs() {
		target := linker.EvalAttrUsage(attr)
		if target&linker.attrTarget["Script"] != 0 {
			contract.RemoveAttr(attr)
			contract.Script().AddAttr(attr)
			continue
		}

		if target&linker.attrTarget["Package"] != 0 {
			contract.RemoveAttr(attr)
			contract.Package().AddAttr(attr)
			continue
		}
		linker.errorf(
			Pos(attr),
			"attr(%s) can't be used to attribute contract :\n\tsee:%s",
			attr,
			Pos(attr.Type.Ref),
		)
	}

	for _, method := range contract.Methods {
		method.Accept(linker)
	}

	return contract
}

//VisitMethod implement visitor interface
func (linker *attrLinker) VisitMethod(method *ast.Method) ast.Node {

	for _, attr := range method.Attrs() {

		target := linker.EvalAttrUsage(attr)

		if target&linker.attrTarget["EnumVal"] == 0 {
			linker.errorf(
				Pos(attr),
				"attr(%s) can't be used to attribute method :\n\tsee:%s",
				attr,
				Pos(attr.Type.Ref),
			)
		}
	}

	for _, expr := range method.Return {
		for _, attr := range expr.Attrs() {
			target := linker.EvalAttrUsage(attr)

			if target&linker.attrTarget["Return"] == 0 {
				linker.errorf(
					Pos(attr),
					"attr(%s) can't be used to attribute method return param :\n\tsee:%s",
					attr,
					Pos(attr.Type.Ref),
				)
			}
		}
	}

	for _, expr := range method.Params {
		for _, attr := range expr.Attrs() {
			target := linker.EvalAttrUsage(attr)

			if target&linker.attrTarget["Param"] == 0 {
				linker.errorf(
					Pos(attr),
					"attr(%s) can't be used to attribute method param :\n\tsee:%s",
					attr,
					Pos(attr.Type.Ref),
				)
			}
		}
	}

	return method
}
