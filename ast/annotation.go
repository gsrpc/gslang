package ast

// Annotation .
type Annotation struct {
	_Node          // mixin _node
	Type  *TypeRef // annotation type
	Args  Expr     // annotation args
}

// NewAnnotation create new comment
func NewAnnotation(name string) *Annotation {
	annotation := &Annotation{
		Type: NewTypeRef(name),
	}

	annotation._init(name)

	return annotation
}
