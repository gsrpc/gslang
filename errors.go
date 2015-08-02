package gslang

import "errors"

// lexer error
var (
	ErrParser = errors.New("gslang parse error")

	ErrDuplicateType = errors.New("duplicate type defined")
)
