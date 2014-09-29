package gslang

import (
	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslang/ast"
)

type evalArg struct {
	field *ast.Field
	expr  ast.Expr
}

//VisitString implement visitor interface
func (visitor *evalArg) VisitString(node *ast.String) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitFloat implement visitor interface
func (visitor *evalArg) VisitFloat(node *ast.Float) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitInt implement visitor interface
func (visitor *evalArg) VisitInt(node *ast.Int) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitBool implement visitor interface
func (visitor *evalArg) VisitBool(node *ast.Bool) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitPackage implement visitor interface
func (visitor *evalArg) VisitPackage(node *ast.Package) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitScript implement visitor interface
func (visitor *evalArg) VisitScript(node *ast.Script) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitEnum implement visitor interface
func (visitor *evalArg) VisitEnum(node *ast.Enum) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitEnumVal implement visitor interface
func (visitor *evalArg) VisitEnumVal(node *ast.EnumVal) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitTable implement visitor interface
func (visitor *evalArg) VisitTable(node *ast.Table) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitField implement visitor interface
func (visitor *evalArg) VisitField(node *ast.Field) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitContract implement visitor interface
func (visitor *evalArg) VisitContract(node *ast.Contract) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitMethod implement visitor interface
func (visitor *evalArg) VisitMethod(node *ast.Method) ast.Node {

	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitTypeRef implement visitor interface
func (visitor *evalArg) VisitTypeRef(node *ast.TypeRef) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitAttr implement visitor interface
func (visitor *evalArg) VisitAttr(node *ast.Attr) ast.Node {

	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitArray implement visitor interface
func (visitor *evalArg) VisitArray(node *ast.Array) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitList implement visitor interface
func (visitor *evalArg) VisitList(node *ast.List) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitArgs implement visitor interface
func (visitor *evalArg) VisitArgs(node *ast.Args) ast.Node {

	for idx, arg := range node.Items {
		if uint16(idx) == visitor.field.ID {
			visitor.expr = arg
		}
	}

	return nil
}

//VisitNamedArgs implement visitor interface
func (visitor *evalArg) VisitNamedArgs(node *ast.NamedArgs) ast.Node {
	for idx, arg := range node.Items {
		if idx == visitor.field.Name() {
			visitor.expr = arg
		}
	}

	return nil
}

//VisitParam implement visitor interface
func (visitor *evalArg) VisitParam(node *ast.Param) ast.Node {
	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil
}

//VisitBinaryOp implement visitor interface
func (visitor *evalArg) VisitBinaryOp(node *ast.BinaryOp) ast.Node {

	gserrors.Panicf(ErrCompileS, "inner error,stmt is not argument list :%s", Pos(node))
	return nil

}
