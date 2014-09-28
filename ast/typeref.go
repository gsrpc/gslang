package ast

import "bytes"

//TypeRef type reference node
type TypeRef struct {
	BasicExpr          //Mixin basic expr implement
	Ref       Expr     //reference node,maybe : contract,enum table
	NamePath  []string //the type reference name path
}

//NewTypeRef create new type reference node
func (node *Script) NewTypeRef(namePath []string) *TypeRef {
	expr := &TypeRef{NamePath: namePath}

	var buff bytes.Buffer
	for _, nodeName := range namePath {
		buff.WriteRune('.')
		buff.WriteString(nodeName)
	}

	expr.Init(buff.String(), node)

	return expr
}
