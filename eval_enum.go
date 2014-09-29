package gslang

import (
	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslang/ast"
)

type evalEnumVal struct {
	val int64
}

//VisitString implement visitor interface
func (visitor *evalEnumVal) VisitString(node *ast.String) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitFloat implement visitor interface
func (visitor *evalEnumVal) VisitFloat(node *ast.Float) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitInt implement visitor interface
func (visitor *evalEnumVal) VisitInt(node *ast.Int) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitBool implement visitor interface
func (visitor *evalEnumVal) VisitBool(node *ast.Bool) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitPackage implement visitor interface
func (visitor *evalEnumVal) VisitPackage(node *ast.Package) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitScript implement visitor interface
func (visitor *evalEnumVal) VisitScript(node *ast.Script) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitEnum implement visitor interface
func (visitor *evalEnumVal) VisitEnum(node *ast.Enum) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitEnumVal implement visitor interface
func (visitor *evalEnumVal) VisitEnumVal(node *ast.EnumVal) ast.Node {
	visitor.val = node.Value
	return node
}

//VisitTable implement visitor interface
func (visitor *evalEnumVal) VisitTable(node *ast.Table) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitField implement visitor interface
func (visitor *evalEnumVal) VisitField(node *ast.Field) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitContract implement visitor interface
func (visitor *evalEnumVal) VisitContract(node *ast.Contract) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitMethod implement visitor interface
func (visitor *evalEnumVal) VisitMethod(node *ast.Method) ast.Node {

	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitTypeRef implement visitor interface
func (visitor *evalEnumVal) VisitTypeRef(node *ast.TypeRef) ast.Node {
	node.Ref.Accept(visitor)
	return node
}

//VisitAttr implement visitor interface
func (visitor *evalEnumVal) VisitAttr(node *ast.Attr) ast.Node {

	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitArray implement visitor interface
func (visitor *evalEnumVal) VisitArray(node *ast.Array) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitList implement visitor interface
func (visitor *evalEnumVal) VisitList(node *ast.List) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitArgs implement visitor interface
func (visitor *evalEnumVal) VisitArgs(node *ast.Args) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitNamedArgs implement visitor interface
func (visitor *evalEnumVal) VisitNamedArgs(node *ast.NamedArgs) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitParam implement visitor interface
func (visitor *evalEnumVal) VisitParam(node *ast.Param) ast.Node {
	gserrors.Panicf(ErrCompileS, "stmt is not const expr :%s", Pos(node))
	return nil
}

//VisitBinaryOp implement visitor interface
func (visitor *evalEnumVal) VisitBinaryOp(node *ast.BinaryOp) ast.Node {

	visitor.val = EvalEnumVal(node.Left) | EvalEnumVal(node.Right)

	return nil
}
