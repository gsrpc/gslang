package ast

// Expr instruction
type Expr interface {
	Node // Mixin Node
}

// ArgsTable .
type ArgsTable struct {
	_Node        // Mixin default node implement
	args  []Expr // arg list
}

// NewArgsTable .
func NewArgsTable() *ArgsTable {
	args := &ArgsTable{}

	args._init("ArgsTable")

	return args
}

// Append .
func (args *ArgsTable) Append(expr Expr) {
	args.args = append(args.args, expr)
}

// Arg .
func (args *ArgsTable) Arg(index int) Expr {
	return args.args[index]
}

// Args .
func (args *ArgsTable) Args() int {
	return len(args.args)
}
