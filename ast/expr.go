package ast

import (
	"errors"
	"fmt"
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
func NewNamedArg(name string) *NamedArg {
	args := &NamedArg{}

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

// Args .
func (args *ArgsTable) Args() int {
	return len(args.args)
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

// Number literal number
type Number struct {
	_Node // Mixin default node implement
	Val   float64
}

// NewNumber .
func NewNumber(val float64) *Number {
	lit := &Number{
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
	_Node // Mixin default node implement
}

// NewConstantRef .
func NewConstantRef(name string) *ConstantRef {
	lit := &ConstantRef{}

	lit._init(name)

	return lit
}

// Op .
type Op int

// Op list
const (
	OpBitwiseOr Op = iota
	OpBitwiseAnd
)

// UnaryOp .
type UnaryOp struct {
	_Node    // Mixin default node implement
	Code  Op // Opcode
}

// NewUnaryOp .
func NewUnaryOp(code Op) *UnaryOp {
	lit := &UnaryOp{
		Code: code,
	}

	lit._init(fmt.Sprintf(""))

	return lit
}
