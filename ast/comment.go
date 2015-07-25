package ast

import "bytes"

// Comment .
type Comment struct {
	_Node              // mixin _node
	buff  bytes.Buffer // comment buff
}

// NewComment create new comment
func NewComment() *Comment {
	comment := &Comment{}

	comment._init("comment")

	return comment
}

// Append .
func (comment *Comment) Append(text string) {
	comment.buff.WriteString(text)
}

func (comment *Comment) String() string {
	return comment.buff.String()
}
