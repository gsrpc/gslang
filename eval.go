package gslang

import (
	"fmt"

	"github.com/gsrpc/gslang/ast"
	"github.com/gsrpc/gslang/lexer"
	"github.com/gsdocker/gslogger"
)

// Eval compile time eval
type Eval interface {
	EvalInt(expr ast.Expr) int64
	EvalString(expr ast.Expr) string
	EvalEnumConstant(name string, constant string) int32
	GetType(name string) (ast.Type, bool)
}

type _Eval struct {
	gslogger.Log              //Mixin logger
	module       *ast.Module  // module
	errorHandler ErrorHandler // error handlers
}

func newEval(errorHandler ErrorHandler, module *ast.Module) Eval {
	return &_Eval{
		Log:          gslogger.Get("eval"),
		module:       module,
		errorHandler: errorHandler,
	}
}

func (eval *_Eval) errorf(err error, node ast.Node, fmtstr string, args ...interface{}) {

	var start lexer.Position
	var end lexer.Position

	if node != nil {
		start, end = Pos(node)
	}

	errinfo := &Error{
		Stage:   StageSemParing,
		Orignal: err,
		Start:   start,
		End:     end,
		Text:    fmt.Sprintf(fmtstr, args...),
	}

	eval.errorHandler.HandleError(errinfo)
}

func (eval *_Eval) EvalString(expr ast.Expr) string {
	stringConstant, ok := expr.(*ast.String)

	if !ok {
		eval.errorf(ErrEval, expr, "can't eval expr(%s) as string ", expr)
	}

	return stringConstant.Name()
}

func (eval *_Eval) EvalEnumConstant(name string, constantName string) int32 {
	target, ok := eval.GetType(name)

	if !ok {
		eval.errorf(ErrEval, nil, "unknown type(%s)", name)
		return 0
	}

	enum, ok := target.(*ast.Enum)

	if !ok {
		eval.errorf(ErrEval, nil, "target type(%s) is not enum", name)
		return 0
	}

	constant, ok := enum.Constant(constantName)

	if !ok {
		eval.errorf(ErrEval, nil, "enum constant %s#%s -- not found", enum, constantName)
		return 0
	}

	return constant.Value
}

func (eval *_Eval) EvalInt(expr ast.Expr) int64 {

	switch expr.(type) {
	case *ast.ConstantRef:
		return eval.EvalInt(expr.(*ast.ConstantRef).Value)
	case *ast.EnumConstant:
		return int64(expr.(*ast.EnumConstant).Value)
	case *ast.BinaryOp:
		binary := expr.(*ast.BinaryOp)

		lhs := eval.EvalInt(binary.LHS)

		rhs := eval.EvalInt(binary.RHS)

		switch binary.Token {
		case lexer.OpBitOr:
			return lhs | rhs
		case lexer.OpBitAnd:
			return lhs & rhs
		default:
			eval.errorf(ErrEval, expr, "can't eval expr(%s) as int64 : unsupport op", expr)
			return 0
		}

	default:
		eval.errorf(ErrEval, expr, "can't eval expr(%s) as int64", expr)
		return 0
	}
}

func (eval *_Eval) GetType(name string) (t ast.Type, ok bool) {
	t, ok = eval.module.Types[name]

	return
}
