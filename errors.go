package gslang

import "errors"

// lexer error
var (
	ErrParser = errors.New("gslang parse error")

	ErrDuplicateType = errors.New("duplicate type defined")

	ErrTypeNotFound = errors.New("not found type")

	ErrFieldName = errors.New("unknown table field name")

	ErrVariableName = errors.New("unknown variable name")

	ErrNewObj = errors.New("newobj ir error")

	ErrAnnotation = errors.New("illegal annotation type")

	ErrType = errors.New("illegal type")

	ErrEval = errors.New("compile time eval error")
)
