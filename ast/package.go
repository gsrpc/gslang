package ast

// Package .
type Package struct {
	_Node                      // Mixin _Node
	scripts map[string]*Script // package contain scripts
}

// NewPackage create new package
func NewPackage(name string) *Package {
	return &Package{
		_Node:   _Node{name: name},
		scripts: make(map[string]*Script),
	}
}

// AddScript .
func (pkg *Package) AddScript(name string, script *Script) {
	pkg.scripts[name] = script
}
