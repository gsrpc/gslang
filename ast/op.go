package ast

//BinaryOp AST binary operation node
type BinaryOp struct {
	BasicExpr      //Mixin basic expr implement
	Left      Expr //Left hand expr
	Right     Expr //right hand expr
}

//NewBinaryOp create new binaryOp
func (node *Script) NewBinaryOp(name string, left Expr, right Expr) *BinaryOp {
	op := &BinaryOp{Left: left, Right: right}
	op.Init(name, node)
	return op
}

//UnaryOp AST binary operation node
type UnaryOp struct {
	BasicExpr      //Mixin basic expr implement
	Right     Expr //Left hand expr
}

//NewUnaryOp create new binaryOp
func (node *Script) NewUnaryOp(name string) *UnaryOp {
	op := &UnaryOp{}
	op.Init(name, node)
	return op
}
