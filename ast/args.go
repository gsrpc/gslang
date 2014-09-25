package ast

//Args AST args node
type Args struct {
	BasicExpr        //Mixin basic expr implement
	Items     []Expr //args list
}

//NewArgs create new args node
func (node *Script) NewArgs() *Args {
	expr := &Args{}
	expr.Init("args", node)
	return expr
}

//NewArg add new arg node to args
func (expr *Args) NewArg(arg Expr) {
	expr.Items = append(expr.Items, arg)
	arg.SetParent(expr)
}

//NamedArgs AST named args node
type NamedArgs struct {
	BasicExpr                 //Mixin basic expr implement
	Items     map[string]Expr //args list
}

//NewNamedArgs create new args node
func (node *Script) NewNamedArgs() *NamedArgs {
	expr := &NamedArgs{Items: make(map[string]Expr)}
	expr.Init("args", node)
	return expr
}

//NewArg add new arg node to args
func (expr *NamedArgs) NewArg(name string, arg Expr) (Expr, bool) {
	if arg, ok := expr.Items[name]; ok {
		return arg, false
	}

	expr.Items[name] = arg
	arg.SetParent(expr)
	return arg, true
}
