package ast

//Package AST's package node
type Package struct {
	BasicNode                    //inher from BasicNode
	Scripts   map[string]*Script //scripts belong to this package
	Types     map[string]Expr    //types belong to this package
}

//NewPackage create new package object
func NewPackage(name string) *Package {
	node := &Package{
		Scripts: make(map[string]*Script),
		Types:   make(map[string]Expr),
	}

	node.Init(name, nil)

	return node
}

//NewType create new type belong to this package
func (node *Package) NewType(expr Expr) (Expr, bool) {
	if old, ok := node.Types[expr.Name()]; ok {
		return old, false
	}

	node.Types[expr.Name()] = expr

	return nil, true
}
