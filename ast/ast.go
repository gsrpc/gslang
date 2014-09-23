package ast

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/gsdocker/gserrors"
)

//Parser parser interface
type Parser interface {
}

//Node ast Node object
type Node interface {
	//Stringer Inher Stringer interface
	fmt.Stringer
	//Name Get the node's display name
	Name() string
	//Path the nodes's full path name in the ast
	Path() string
	//Parent get this node's parent node,may be nil
	Parent() Node
	//SetParent set this node's parent ,return old parent node
	SetParent(parent Node) Node
	//Package The node belongs Package
	Package() *Package
	//Parser The node belongs parser
	Parser() Parser
	//Attrs get the attrs table
	Attrs() []*Attr
	//DelAttr Delete target attr from attribute table
	DelAttr(attr *Attr)
	//NewAttr add new attr to attribute table
	NewAttr(attr *Attr)
	//NewExtra add new extra data,return old extra data with same name
	NewExtra(name string, data interface{})
	//Extra get extra data by name
	Extra(name string) (interface{}, bool)
	//DelExtra delete extra data by name
	DelExtra(name string)
	//Accept implement visitor design pattern
	Accept(visitor Visitor)
}

//Path Get node's AST path
func Path(node Node) (result []Node) {
	var nodes []Node

	current := node

	for current != nil {
		nodes = append(nodes, current)
		current = current.Parent()
	}

	for i := len(nodes) - 1; i > -1; i-- {
		result = append(result, nodes[i])
	}

	return
}

//Expr AST expr node
type Expr interface {
	//Node Inher Node interface
	Node
	//Script The type belongs script
	Script() *Script
}

//BasicNode node default implement
type BasicNode struct {
	name   string                 //名称
	parent Node                   //父节点
	attrs  []*Attr                //属性列表
	extras map[string]interface{} //附加数据表
}

//Init implement Node interface
func (node *BasicNode) Init(name string, parent Node) {
	node.name = name
	node.parent = parent
}

//Name implement Node interface
func (node *BasicNode) Name() string {
	return node.name
}

//String implement Stringer interface
func (node *BasicNode) String() string {
	return node.name
}

//Path implement Node interface
func (node *BasicNode) Path() string {
	var writer bytes.Buffer

	for _, node := range Path(node) {
		writer.WriteString(node.Name())
		writer.WriteRune('.')
	}

	return writer.String()
}

//Parser implement Node interface
func (node *BasicNode) Parser() Parser {
	if node.Parent() == nil {
		return nil
	}

	return node.Parent().Parser()
}

//Package implement Node interface
func (node *BasicNode) Package() *Package {
	if node.Parent() == nil {
		return nil
	}

	return node.Parent().Package()
}

//Parent implement Node interface
func (node *BasicNode) Parent() Node {
	return node.parent
}

//SetParent implement Node interface
func (node *BasicNode) SetParent(parent Node) (old Node) {
	old, node.parent = node.parent, parent
	return
}

func (node *BasicNode) getExtra() map[string]interface{} {
	if node.extras == nil {
		node.extras = make(map[string]interface{})
	}

	return node.extras
}

//Attrs implement Node interface
func (node *BasicNode) Attrs() []*Attr {
	return node.attrs
}

//NewAttr implement Node interface
func (node *BasicNode) NewAttr(attr *Attr) {
	for _, old := range node.attrs {
		if old == attr {
			return
		}
	}

	attr.SetParent(node)

	node.attrs = append(node.attrs, attr)
}

//DelAttr implement Node interface
func (node *BasicNode) DelAttr(attr *Attr) {
	var attrs []*Attr
	for _, old := range node.attrs {
		if old == attr {

			continue
		}
		attrs = append(attrs, old)
	}

	attr.SetParent(nil)
}

//NewExtra implement Node interface
func (node *BasicNode) NewExtra(name string, data interface{}) {
	node.getExtra()[name] = data
}

//Extra implement Node interface
func (node *BasicNode) Extra(name string) (data interface{}, ok bool) {
	data, ok = node.getExtra()[name]
	return
}

//DelExtra implement Node interface
func (node *BasicNode) DelExtra(name string) {
	delete(node.getExtra(), name)
}

//Accept implement Node interface
func (node *BasicNode) Accept(Visitor) {
	gserrors.Panicf(nil, "type(%s) not implement Accept", reflect.TypeOf(node))
}

//BasicExpr default implement of Expr interface
type BasicExpr struct {
	BasicNode
	script *Script
}

//Init implement Expr interface
func (node *BasicExpr) Init(name string, script *Script) {
	defer gserrors.Ensure(script != nil, "the script param can't be nil")
	node.BasicNode.Init(name, nil)
	node.script = script
}

//Script implement Expr interface
func (node *BasicExpr) Script() *Script {
	gserrors.Require(node.script != nil, "the script param can't be nil")
	return node.script
}