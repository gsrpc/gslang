package gslang

import (
	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslang/ast"
)

func (cs *CompileS) link(pkg *ast.Package) {
	linker := &Linker{
		CompileS: cs,
	}

	pkg.Accept(linker)
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

	for _, attr := range script.Attrs() {
		linker.EvalAttrUsage(attr)
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
