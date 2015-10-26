package gslang

import (
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

// FindAnnotations .
func FindAnnotations(node ast.Node, name string) (retval []*ast.Annotation) {

	val, ok := node.GetExtra(ExtraAnnotation)

	if ok {
		anns := val.([]*ast.Annotation)

		for _, ann := range anns {

			if ann.Type.Ref == nil {
				continue
			}

			if ann.Type.Ref.FullName() == name {
				retval = append(retval, ann)
			}
		}
	}

	return
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

// EnumSize get enum binary size
func EnumSize(typeDecl ast.Type) int {
	_, ok := FindAnnotation(typeDecl, "gslang.Flag")

	if ok {
		return 4
	}

	return 1
}

// EnumType get enum token type
func EnumType(typeDecl ast.Type) lexer.TokenType {
	_, ok := FindAnnotation(typeDecl, "gslang.Flag")

	if ok {
		return lexer.KeyUInt32
	}

	return lexer.KeyByte
}

// IsBuiltin check if target type is builtin type
func IsBuiltin(typeDecl ast.Type) bool {
	_, ok := typeDecl.(*ast.BuiltinType)

	return ok
}

// IsVoid check if target type is builtin type Void
func IsVoid(typeDecl ast.Type) bool {
	builtinType, ok := typeDecl.(*ast.BuiltinType)

	if ok && builtinType.Type == lexer.KeyVoid {
		return true
	}

	return false
}

// NotVoid check if target type is not builtin type Void
func NotVoid(typeDecl ast.Type) bool {
	return !IsVoid(typeDecl)
}

// IsPOD check if target type is POD table
func IsPOD(typeDecl ast.Type) bool {
	_, ok := FindAnnotation(typeDecl, "gslang.POD")

	if ok {
		return true
	}

	return false
}

// IsAsync check if target method is async method
func IsAsync(method *ast.Method) bool {
	_, ok := FindAnnotation(method, "gslang.Async")

	if ok {
		return true
	}

	return false
}

// IsException check if target type is exception table
func IsException(typeDecl ast.Type) bool {
	_, ok := FindAnnotation(typeDecl, "gslang.Exception")

	if ok {
		return true
	}

	return false
}
