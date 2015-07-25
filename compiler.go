package gslang

import (
	"bytes"
	"io/ioutil"

	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslogger"
)

// Extra key list
const (
	ExtraStartPos = "start"
	ExtraEndPos   = "end"
	ExtraComment  = "comment"
)

func _setNodePos(node ast.Node, start Position, end Position) {
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

// Pos .
func Pos(node ast.Node) (start Position, end Position) {
	node.GetExtra(ExtraStartPos, &start)
	node.GetExtra(ExtraEndPos, &end)

	return
}

// ErrorHandler .
type ErrorHandler interface {
	// handle
	HandleParseError(err error, position Position, msg string)
}

// HandleParseError .
type HandleParseError func(err error, position Position, msg string)

// HandleParseError implement ErrorHandler
func (handle HandleParseError) HandleParseError(err error, position Position, msg string) {
	handle(err, position, msg)
}

// Compiler gslang compiler
type Compiler struct {
	gslogger.Log                        // Mixin log
	scripts      map[string]*ast.Script // compiled scripts
}

// NewCompiler .
func NewCompiler() *Compiler {
	return &Compiler{
		Log:     gslogger.Get("compiler"),
		scripts: make(map[string]*ast.Script),
	}
}

// Compile .
func (compiler *Compiler) Compile(filepath string, errorHandler ErrorHandler) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	content, err := ioutil.ReadFile(filepath)

	if err != nil {
		return err
	}

	compiler.scripts[filepath] = compiler.parse(NewLexer(filepath, bytes.NewBuffer(content)), errorHandler)

	return
}
