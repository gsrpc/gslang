package ast

import "fmt"

// Node .
type Node interface {
	fmt.Stringer                             // Mixin Stringer
	Name() string                            // node name
	Parent() Node                            // parent node
	SetParent(node Node) Node                // set parent node
	SetExtra(key string, val interface{})    // set extra data
	GetExtra(key string) (interface{}, bool) // get extra data
}

type _Node struct {
	name   string                 // node name
	parent Node                   // parent
	extra  map[string]interface{} // extra data map
}

func (node *_Node) _init(name string) {
	node.name = name
	node.extra = make(map[string]interface{})
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

func (node *_Node) SetExtra(key string, val interface{}) {
	node.extra[key] = val
}

func (node *_Node) GetExtra(key string) (val interface{}, ok bool) {
	val, ok = node.extra[key]

	return
}
