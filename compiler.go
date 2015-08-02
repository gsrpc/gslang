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

func _AttachAnnotation(node ast.Node, annotation []*ast.Annotation) {

	anns := Annotation(node)

	anns = append(anns, annotation...)

	node.SetExtra(ExtraComment, anns)
}

// Annotation .
func Annotation(node ast.Node) (anns []*ast.Annotation) {
	node.GetExtra(ExtraAnnotation, &anns)

	return
}

// Pos .
func Pos(node ast.Node) (start lexer.Position, end lexer.Position) {
	node.GetExtra(ExtraStartPos, &start)
	node.GetExtra(ExtraEndPos, &end)

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
	gslogger.Log                        // Mixin log
	scripts      map[string]*ast.Script // compiled scripts
	errorHandler ErrorHandler
}

// NewCompiler .
func NewCompiler(errorHandler ErrorHandler) *Compiler {
	return &Compiler{
		Log:          gslogger.Get("compiler"),
		scripts:      make(map[string]*ast.Script),
		errorHandler: errorHandler,
	}
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

	compiler.scripts[filepath] = compiler.parse(lexer.NewLexer(filepath, bytes.NewBuffer(content)), compiler.errorHandler)

	return
}
