package gslang

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslogger"
	"github.com/gsrpc/gslang/ast"
	"github.com/gsrpc/gslang/lexer"
)

// Extra key list
const (
	ExtraStartPos   = "start"
	ExtraEndPos     = "end"
	ExtraComment    = "comment"
	ExtraAnnotation = "annotation"
)

func _setNodePos(node ast.Node, start lexer.Position, end lexer.Position) {
	node.SetExtra(ExtraStartPos, start)
	node.SetExtra(ExtraEndPos, end)
}

func _AttachComment(node ast.Node, comment *ast.Comment) bool {
	nodeStart, nodeEnd := Pos(node)

	commentStart, commentEnd := Pos(comment)

	if nodeStart.Lines == commentEnd.Lines ||
		nodeStart.Lines == commentEnd.Lines+1 ||
		nodeEnd.Lines == commentStart.Lines {

		node.SetExtra(ExtraComment, comment)

		return true
	}

	return false
}

func _AttachAnnotation(node ast.Node, annotations ...*ast.Annotation) {

	anns := Annotations(node)

	anns = append(anns, annotations...)

	node.SetExtra(ExtraAnnotation, anns)
}

// Annotations .
func Annotations(node ast.Node) (anns []*ast.Annotation) {

	val, ok := node.GetExtra(ExtraAnnotation)

	if ok {
		anns = val.([]*ast.Annotation)
	}

	return
}

// Annotations .
func _RemoveAnnotation(node ast.Node, annotation *ast.Annotation) {

	anns := Annotations(node)

	var newanns []*ast.Annotation

	for _, ann := range anns {
		if ann == annotation {
			continue
		}
		newanns = append(newanns, ann)
	}
	node.SetExtra(ExtraAnnotation, anns)
}

// FindAnnotation .
func FindAnnotation(node ast.Node, name string) (*ast.Annotation, bool) {

	val, ok := node.GetExtra(ExtraAnnotation)

	if ok {
		anns := val.([]*ast.Annotation)

		for _, ann := range anns {

			if ann.Type.Ref == nil {
				continue
			}

			if ann.Type.Ref.FullName() == name {
				return ann, true
			}
		}
	}

	return nil, false
}

// Pos .
func Pos(node ast.Node) (start lexer.Position, end lexer.Position) {
	val, ok := node.GetExtra(ExtraStartPos)

	if ok {
		start = val.(lexer.Position)
	}

	val, ok = node.GetExtra(ExtraEndPos)

	if ok {
		end = val.(lexer.Position)
	}

	return
}

// Stage compile stage
type Stage int

// compile stages
const (
	StageLexer Stage = iota
	StageParing
	StageSemParing
)

// Error compile error context
type Error struct {
	Stage   Stage          // compile stage
	Orignal error          // orignal error code
	Start   lexer.Position // error location start
	End     lexer.Position // error location end
	Text    string         // error description
}

// ErrorHandler .
type ErrorHandler interface {
	HandleError(err *Error)
}

// HandleError .
type HandleError func(err *Error)

// HandleError implement ErrorHandler
func (handle HandleError) HandleError(err *Error) {
	handle(err)
}

// Compiler gslang compiler
type Compiler struct {
	gslogger.Log              // Mixin log
	module       *ast.Module  // compiled scripts
	errorHandler ErrorHandler // error handler
	eval         Eval         //eval site
}

// NewCompiler .
func NewCompiler(name string, errorHandler ErrorHandler) *Compiler {

	module := ast.NewModule(name)

	return &Compiler{
		Log:          gslogger.Get("compiler"),
		module:       module,
		errorHandler: errorHandler,
		eval:         newEval(errorHandler, module),
	}

}

// Eval .
func (compiler *Compiler) Eval() Eval {
	return compiler.eval
}

// Compile .
func (compiler *Compiler) Compile(filepath string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	content, err := ioutil.ReadFile(filepath)

	if err != nil {
		return err
	}

	compiler.parse(lexer.NewLexer(filepath, bytes.NewBuffer(content)), compiler.errorHandler)

	return
}

// Visitor gslang CodeGen
type Visitor interface {
	BeginScript(compiler *Compiler, script *ast.Script) bool
	// get using template
	Using(compiler *Compiler, using *ast.Using)

	Table(compiler *Compiler, tableType *ast.Table)

	Exception(compiler *Compiler, tableType *ast.Table)

	Annotation(compiler *Compiler, annotation *ast.Table)

	Enum(compiler *Compiler, enum *ast.Enum)

	Contract(compiler *Compiler, contract *ast.Contract)
	//
	EndScript(compiler *Compiler)
}

type _Visitor struct {
	gslogger.Log             // log APIs
	codeGen      Visitor     //implement
	module       *ast.Module // generate code module
	compiler     *Compiler
}

// Visit visit ast tree
func (compiler *Compiler) Visit(codeGen Visitor) (err error) {

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

	gen := &_Visitor{
		Log:      gslogger.Get("codegen"),
		codeGen:  codeGen,
		module:   compiler.module,
		compiler: compiler,
	}

	gen.visit()

	return
}

func (codeGen *_Visitor) visit() {

	codeGen.module.Foreach(func(script *ast.Script) bool {

		if !codeGen.codeGen.BeginScript(codeGen.compiler, script) {
			return true
		}

		script.UsingForeach(func(using *ast.Using) {

			if strings.HasPrefix(using.Name(), "gslang") || strings.HasPrefix(using.Name(), "gslang.annotations") {
				return
			}

			codeGen.codeGen.Using(codeGen.compiler, using)
		})

		script.TypeForeach(func(typeDecl ast.Type) {
			switch typeDecl.(type) {
			case *ast.Table:
				_, ok := FindAnnotation(typeDecl, "gslang.annotations.Usage")

				if ok {
					codeGen.codeGen.Annotation(codeGen.compiler, typeDecl.(*ast.Table))
					break
				}

				_, ok = FindAnnotation(typeDecl, "gslang.Exception")

				if ok {
					codeGen.codeGen.Exception(codeGen.compiler, typeDecl.(*ast.Table))
					break
				}

				codeGen.codeGen.Table(codeGen.compiler, typeDecl.(*ast.Table))

			case *ast.Enum:
				codeGen.codeGen.Enum(codeGen.compiler, typeDecl.(*ast.Enum))
			case *ast.Contract:
				codeGen.codeGen.Contract(codeGen.compiler, typeDecl.(*ast.Contract))
			}
		})

		codeGen.codeGen.EndScript(codeGen.compiler)

		return true
	})
}
