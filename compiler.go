package gslang

import (
	"bytes"
	"io/ioutil"

	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslang/lexer"
	"github.com/gsdocker/gslogger"
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

func _AttachAnnotation(node ast.Node, annotations []*ast.Annotation) {

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
		eval:         newEval(module),
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
