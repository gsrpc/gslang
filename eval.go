package gslang

import (
	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslogger"
)

// Eval compile time eval
type Eval interface {
	EvalInt(expr ast.Expr) int64
	GetType(name string) (ast.Type,bool)
}

type _Eval struct {
	gslogger.Log             //Mixin logger
	module       *ast.Module // module
}

func newEval(module *ast.Module) Eval {
	return &_Eval{
		Log:    gslogger.Get("eval"),
		module: module,
	}
}

func (eval *_Eval) EvalInt(expr ast.Expr) int64 {
	return 0
}

func (eval *_Eval) GetType(name string) (t ast.Type, ok bool) {
	t, ok = eval.module.Types[name]

	return
}
