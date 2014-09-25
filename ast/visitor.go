package ast

//Visitor visitor interface
type Visitor interface {
	VisitPackage(*Package)
	VisitScript(*Script)
	VisitEnum(*Enum)
	VisitEnumVal(*EnumVal)
	VisitTable(*Table)
	VisitField(*Field)
	VisitContract(*Contract)
	VisitMethod(*Method)
	VisitTypeRef(*TypeRef)
	VisitAttr(*Attr)
	VisitArray(*Array)
}

//Accept implement Node interface
func (node *Package) Accept(visitor Visitor) {
	visitor.VisitPackage(node)
}

//Accept implement Node interface
func (node *Script) Accept(visitor Visitor) {
	visitor.VisitScript(node)
}

//Accept implement Node interface
func (node *Enum) Accept(visitor Visitor) {
	visitor.VisitEnum(node)
}

//Accept implement Node interface
func (node *EnumVal) Accept(visitor Visitor) {
	visitor.VisitEnumVal(node)
}

//Accept implement Node interface
func (node *Table) Accept(visitor Visitor) {
	visitor.VisitTable(node)
}

//Accept implement Node interface
func (node *Field) Accept(visitor Visitor) {
	visitor.VisitField(node)
}

//Accept implement Node interface
func (node *Contract) Accept(visitor Visitor) {
	visitor.VisitContract(node)
}

//Accept implement Node interface
func (node *Method) Accept(visitor Visitor) {
	visitor.VisitMethod(node)
}

//Accept implement Node interface
func (node *TypeRef) Accept(visitor Visitor) {
	visitor.VisitTypeRef(node)
}

//Accept implement Node interface
func (node *Attr) Accept(visitor Visitor) {
	visitor.VisitAttr(node)
}

//Accept implement Node interface
func (node *Array) Accept(visitor Visitor) {
	visitor.VisitArray(node)
}
