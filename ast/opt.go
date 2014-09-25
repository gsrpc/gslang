package ast

//BinaryOp AST binary operation node
type BinaryOp struct {
	BasicExpr      //Mixin basic expr implement
	Left      Expr //Left hand expr
	Right     Expr //right hand expr
}

//NewBinaryOp create new binaryOp
func (node *Script) NewBinaryOp(name string) *BinaryOp {
	op := &BinaryOp{}
	op.Init(name, node)
	return op
}
