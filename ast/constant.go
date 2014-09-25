package ast

//String literal string node
type String struct {
	BasicExpr        //Mixin basic expr implement
	Value     string //literal string content
}

//NewString create new literal string node
func (node *Script) NewString(content string) *String {
	expr := &String{Value: content}
	expr.Init("string", node)
	return expr
}

//Float literal float node
type Float struct {
	BasicExpr         //Mixin basic expr implement
	Value     float64 //literal string content
}

//NewFloat create new literal string node
func (node *Script) NewFloat(val float64) *Float {
	expr := &Float{Value: val}
	expr.Init("float", node)
	return expr
}

//Int literal int node
type Int struct {
	BasicExpr       //Mixin basic expr implement
	Value     int64 //literal string content
}

//NewInt create new literal integer node
func (node *Script) NewInt(val int64) *Int {
	expr := &Int{Value: val}
	expr.Init("int", node)
	return expr
}
