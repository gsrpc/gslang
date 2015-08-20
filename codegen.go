package gslang

import (
	"text/template"

	"github.com/gsdocker/gslang/ast"
)

// CodeTarget language code target
type CodeTarget interface {
	Using() *template.Template
	Table() *template.Template
	Exception() *template.Template
	Annotations() *template.Template
	Enum() *template.Template
	Contract() *template.Template

	Begin()
	CreateScript(script *ast.Script, header string)
	End()
}

type _CodeGen struct {
	target CodeTarget
}

// Gen generate code
func (compiler *Compiler) Gen(target CodeTarget) (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	codegen := &_CodeGen{
		target: target,
	}

	return compiler.Visit(codegen)
}

func (codeGen *_CodeGen) BeginScript(script *ast.Script) {

}

func (codeGen *_CodeGen) Using(using *ast.Using) {

}
func (codeGen *_CodeGen) Table(tableType *ast.Table) {

}
func (codeGen *_CodeGen) Exception(tableType *ast.Table) {

}
func (codeGen *_CodeGen) Annotation(annotation *ast.Table) {

}
func (codeGen *_CodeGen) Enum(enum *ast.Enum) {

}
func (codeGen *_CodeGen) Contract(contract *ast.Contract) {

}
func (codeGen *_CodeGen) EndScript() {

}
