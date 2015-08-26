package ast

import (
	"errors"
	"fmt"

	"github.com/gsdocker/gslang/lexer"
)

// errors .
var (
	ErrArgType = errors.New("unexpect arg type")
)

// Expr instruction
type Expr interface {
	Node // Mixin Node
}

// NamedArg .
type NamedArg struct {
	_Node      // Mixin default node implement
	Arg   Expr // named arg
}

// NewNamedArg new named arg
func NewNamedArg(name string, arg Expr) *NamedArg {
	args := &NamedArg{
		Arg: arg,
	}

	args._init(name)

	return args
}

// ArgsTable .
type ArgsTable struct {
	_Node        // Mixin default node implement
	args  []Expr // arg list
	Named bool   // named args table flag
}

// NewArgsTable .
func NewArgsTable(named bool) *ArgsTable {
	args := &ArgsTable{
		Named: named,
	}

	args._init("ArgsTable")

	return args
}

// Append .
func (args *ArgsTable) Append(expr Expr) error {

	_, named := expr.(*NamedArg)

	if args.Named != named {
		return ErrArgType
	}

	args.args = append(args.args, expr)

	return nil
}

// Arg .
func (args *ArgsTable) Arg(index int) Expr {
	return args.args[index]
}

// NamedArg .
func (args *ArgsTable) NamedArg(name string) (Expr, bool) {
	for _, arg := range args.args {
		namedArg := arg.(*NamedArg)
		if namedArg.Name() == name {
			return namedArg.Arg, true
		}
	}

	return nil, false
}

// Count .
func (args *ArgsTable) Count() int {
	return len(args.args)
}

// Args .
func (args *ArgsTable) Args() []Expr {
	return args.args
}

// String literal string
type String struct {
	_Node // Mixin default node implement
}

// NewString .
func NewString(val string) *String {
	lit := &String{}

	lit._init(val)

	return lit
}

// Numeric literal number
type Numeric struct {
	_Node // Mixin default node implement
	Val   float64
}

// NewNumeric .
func NewNumeric(val float64) *Numeric {
	lit := &Numeric{
		Val: val,
	}

	lit._init(fmt.Sprintf("%f", val))

	return lit
}

// Boolean literal boolean
type Boolean struct {
	_Node // Mixin default node implement
	Val   bool
}

// NewBoolean .
func NewBoolean(val bool) *Boolean {
	lit := &Boolean{
		Val: val,
	}

	lit._init(fmt.Sprintf("%t", val))

	return lit
}

// ConstantRef .
type ConstantRef struct {
	_Node      // Mixin default node implement
	Value Expr // constant val
}

// NewConstantRef .
func NewConstantRef(name string) *ConstantRef {
	lit := &ConstantRef{}

	lit._init(name)

	return lit
}

// NewObj .
type NewObj struct {
	_Node            // Mixin default node implement
	Type  *TypeRef   // obj type reference
	Args  *ArgsTable // new initialize list
}

// NewNewObj .
func NewNewObj(name string, args *ArgsTable) *NewObj {
	lit := &NewObj{
		Args: args,
		Type: NewTypeRef(name),
	}

	lit._init(name)

	return lit
}

// NewNewObj2 .
func NewNewObj2(ref *TypeRef, args *ArgsTable) *NewObj {
	lit := &NewObj{
		Args: args,
		Type: ref,
	}

	lit._init(ref.Name())

	return lit
}

// OpType type
type OpType int

// Op type list
const (
	OpBinary OpType = 1 << iota
	OpUnary
)

// UnaryOp .
type UnaryOp struct {
	_Node                   // Mixin default node implement
	Token   lexer.TokenType // op code
	Operand Expr            // Binary op object
}

// NewUnaryOp .
func NewUnaryOp(token lexer.TokenType, operand Expr) *UnaryOp {
	lit := &UnaryOp{
		Token:   token,
		Operand: operand,
	}

	lit._init(token.String())

	return lit
}

// BinaryOp .
type BinaryOp struct {
	_Node                 // Mixin default node implement
	Token lexer.TokenType // op code
	LHS   Expr            // Binary op left handle experand
	RHS   Expr            // Binary op right handle experand
}

// NewBinaryOp .
func NewBinaryOp(token lexer.TokenType, lhs, rhs Expr) *BinaryOp {
	lit := &BinaryOp{
		Token: token,
		LHS:   lhs,
		RHS:   rhs,
	}

	lit._init(token.String())

	return lit
}
