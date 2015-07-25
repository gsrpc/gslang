package gslang

import "errors"

// lexer error
var (
	ErrLexer  = errors.New("gslang lexer error")
	ErrParser = errors.New("gslang parse error")
)
