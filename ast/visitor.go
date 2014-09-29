package ast

import "errors"

//visitor  errors
var (
	ErrVisit = errors.New("invalid call to visit method")
)

//Visitor visitor interface
type Visitor interface {
	VisitPackage(*Package) Node
	VisitScript(*Script) Node
	VisitEnum(*Enum) Node
	VisitEnumVal(*EnumVal) Node
	VisitTable(*Table) Node
	VisitField(*Field) Node
	VisitContract(*Contract) Node
	VisitMethod(*Method) Node
	VisitParam(*Param) Node
	VisitTypeRef(*TypeRef) Node
	VisitAttr(*Attr) Node
	VisitArray(*Array) Node
	VisitList(*List) Node
	VisitArgs(*Args) Node
	VisitNamedArgs(*NamedArgs) Node
	VisitString(*String) Node
	VisitFloat(*Float) Node
	VisitInt(*Int) Node
	VisitBool(*Bool) Node
	VisitBinaryOp(*BinaryOp) Node
}

//Accept implement Node interface
func (node *BinaryOp) Accept(visitor Visitor) Node {
	return visitor.VisitBinaryOp(node)
}

//Accept implement Node interface
func (node *Param) Accept(visitor Visitor) Node {
	return visitor.VisitParam(node)
}

//Accept implement Node interface
func (node *String) Accept(visitor Visitor) Node {
	return visitor.VisitString(node)
}

//Accept implement Node interface
func (node *Float) Accept(visitor Visitor) Node {
	return visitor.VisitFloat(node)
}

//Accept implement Node interface
func (node *Int) Accept(visitor Visitor) Node {
	return visitor.VisitInt(node)
}

//Accept implement Node interface
func (node *Bool) Accept(visitor Visitor) Node {
	return visitor.VisitBool(node)
}

//Accept implement Node interface
func (node *NamedArgs) Accept(visitor Visitor) Node {
	return visitor.VisitNamedArgs(node)
}

//Accept implement Node interface
func (node *Args) Accept(visitor Visitor) Node {
	return visitor.VisitArgs(node)
}

//Accept implement Node interface
func (node *Package) Accept(visitor Visitor) Node {
	return visitor.VisitPackage(node)
}

//Accept implement Node interface
func (node *Script) Accept(visitor Visitor) Node {
	return visitor.VisitScript(node)
}

//Accept implement Node interface
func (node *Enum) Accept(visitor Visitor) Node {
	return visitor.VisitEnum(node)
}

//Accept implement Node interface
func (node *EnumVal) Accept(visitor Visitor) Node {
	return visitor.VisitEnumVal(node)
}

//Accept implement Node interface
func (node *Table) Accept(visitor Visitor) Node {
	return visitor.VisitTable(node)
}

//Accept implement Node interface
func (node *Field) Accept(visitor Visitor) Node {
	return visitor.VisitField(node)
}

//Accept implement Node interface
func (node *Contract) Accept(visitor Visitor) Node {
	return visitor.VisitContract(node)
}

//Accept implement Node interface
func (node *Method) Accept(visitor Visitor) Node {
	return visitor.VisitMethod(node)
}

//Accept implement Node interface
func (node *TypeRef) Accept(visitor Visitor) Node {
	return visitor.VisitTypeRef(node)
}

//Accept implement Node interface
func (node *Attr) Accept(visitor Visitor) Node {
	return visitor.VisitAttr(node)
}

//Accept implement Node interface
func (node *Array) Accept(visitor Visitor) Node {
	return visitor.VisitArray(node)
}

//Accept implement Node interface
func (node *List) Accept(visitor Visitor) Node {
	return visitor.VisitList(node)
}

///////////////////////////////////////////////////

//EmptyVisitor do nothing but throw an assert exception
type EmptyVisitor struct{}

//VisitString implement visitor interface
func (visitor *EmptyVisitor) VisitString(*String) Node {
	return nil
}

//VisitFloat implement visitor interface
func (visitor *EmptyVisitor) VisitFloat(*Float) Node {
	return nil
}

//VisitInt implement visitor interface
func (visitor *EmptyVisitor) VisitInt(*Int) Node {
	return nil
}

//VisitBool implement visitor interface
func (visitor *EmptyVisitor) VisitBool(*Bool) Node {
	return nil
}

//VisitPackage implement visitor interface
func (visitor *EmptyVisitor) VisitPackage(*Package) Node {
	return nil
}

//VisitScript implement visitor interface
func (visitor *EmptyVisitor) VisitScript(*Script) Node {
	return nil
}

//VisitEnum implement visitor interface
func (visitor *EmptyVisitor) VisitEnum(*Enum) Node {
	return nil
}

//VisitEnumVal implement visitor interface
func (visitor *EmptyVisitor) VisitEnumVal(*EnumVal) Node {
	return nil
}

//VisitTable implement visitor interface
func (visitor *EmptyVisitor) VisitTable(*Table) Node {
	return nil
}

//VisitField implement visitor interface
func (visitor *EmptyVisitor) VisitField(*Field) Node {
	return nil
}

//VisitContract implement visitor interface
func (visitor *EmptyVisitor) VisitContract(*Contract) Node {
	return nil
}

//VisitMethod implement visitor interface
func (visitor *EmptyVisitor) VisitMethod(*Method) Node {

	return nil
}

//VisitTypeRef implement visitor interface
func (visitor *EmptyVisitor) VisitTypeRef(*TypeRef) Node {
	return nil
}

//VisitAttr implement visitor interface
func (visitor *EmptyVisitor) VisitAttr(*Attr) Node {

	return nil
}

//VisitArray implement visitor interface
func (visitor *EmptyVisitor) VisitArray(*Array) Node {
	return nil
}

//VisitList implement visitor interface
func (visitor *EmptyVisitor) VisitList(*List) Node {
	return nil
}

//VisitArgs implement visitor interface
func (visitor *EmptyVisitor) VisitArgs(*Args) Node {
	return nil
}

//VisitNamedArgs implement visitor interface
func (visitor *EmptyVisitor) VisitNamedArgs(*NamedArgs) Node {
	return nil
}

//VisitParam implement visitor interface
func (visitor *EmptyVisitor) VisitParam(*Param) Node {
	return nil
}

//VisitBinaryOp implement visitor interface
func (visitor *EmptyVisitor) VisitBinaryOp(*BinaryOp) Node {
	return nil
}

///////////////////////////////////////////////////////////////
