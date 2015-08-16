package gslang

import (
	"text/template"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslogger"
)

// CodeGen gslang CodeGen
type CodeGen interface {
	BeginScript(script *ast.Script)
	// get using template
	Using(using *ast.Using)

	Table(tableType *ast.Table)

	Annotation(annotation *ast.Table)

	Enum(enum *ast.Enum)

	Contract(contract *ast.Contract)
	//
	EndScript()
}

type _CodeGen struct {
	gslogger.Log                    // log APIs
	codeGen      CodeGen            //implement
	module       *ast.Module        // generate code module
	usingTmplate *template.Template // using template
}

// Gen generate codes
func (compiler *Compiler) Gen(codeGen CodeGen) (err error) {

	defer func() {
		if e := recover(); e != nil {

			gserr, ok := e.(gserrors.GSError)

			if ok {
				err = gserr
			} else {
				err = gserrors.Newf(e.(error), "catch unknown error")
			}
		}
	}()

	gen := &_CodeGen{
		Log:     gslogger.Get("codegen"),
		codeGen: codeGen,
		module:  compiler.module,
	}

	gen.Gen()

	return
}

func (codeGen *_CodeGen) Gen() {

	codeGen.module.Foreach(func(script *ast.Script) bool {

		codeGen.codeGen.BeginScript(script)

		script.UsingForeach(func(using *ast.Using) {
			codeGen.codeGen.Using(using)
		})

		script.TypeForeach(func(typeDecl ast.Type) {
			switch typeDecl.(type) {
			case *ast.Table:
				_, ok := FindAnnotation(typeDecl, "gslang.annotations.Usage")

				if ok {
					codeGen.codeGen.Annotation(typeDecl.(*ast.Table))
				} else {
					codeGen.codeGen.Table(typeDecl.(*ast.Table))
				}

			case *ast.Enum:
				codeGen.codeGen.Enum(typeDecl.(*ast.Enum))
			case *ast.Contract:
				codeGen.codeGen.Contract(typeDecl.(*ast.Contract))
			}
		})

		codeGen.codeGen.EndScript()

		return true
	})
}
