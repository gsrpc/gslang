package ast

import "fmt"

// Node .
type Node interface {
	fmt.Stringer              // Mixin Stringer
	Name() string             // node name
	Parent() Node             // parent node
	SetParent(node Node) Node // set parent node
}

type _Node struct {
	name   string // node name
	parent Node   // parent
}

//Name implement Node interface
func (node *_Node) Name() string {
	return node.name
}

//String implement Stringer interface
func (node *_Node) String() string {
	return node.name
}

//Parent implement Node interface
func (node *_Node) Parent() Node {
	return node.parent
}

//SetParent implement Node interface
func (node *_Node) SetParent(parent Node) (old Node) {
	old, node.parent = node.parent, parent
	return
}
