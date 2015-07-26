package ast

import "bytes"

// Annotation .
type Annotation struct {
	_Node              // mixin _node
	buff  bytes.Buffer // comment buff
}

// NewAnnotation create new comment
func NewAnnotation(name string) *Annotation {
	annotation := &Annotation{}

	annotation._init(name)

	return annotation
}
