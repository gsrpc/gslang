package ast

import "github.com/gsdocker/gserrors"

//Script AST script node
type Script struct {
	BasicNode                     //Inher from default Node
	Imports   map[string]*Package //import packages
	Types     []Expr              //types defined in this script
	pkg       *Package            //script belongs package
}

//NewScript create new script node
func (node *Package) NewScript(name string) (script *Script, err error) {
	//check if created script with the same name
	if old, ok := node.Scripts[name]; ok {
		script = old
		err = gserrors.Newf(ErrAST, "duplicate script named :%s", name)
		return
	}

	defer gserrors.Ensure(func() bool {
		return script.pkg != nil
	}, "make sure set the Package field")

	defer gserrors.Ensure(func() bool {
		return script.Imports != nil
	}, "make sure init the Imports field")

	script = &Script{
		Imports: make(map[string]*Package),
		pkg:     node,
	}
	//call BasicNode's init function
	script.Init(name, node)

	node.Scripts[name] = script

	return
}

//NewType create new type node for this script
//if the return param ok is false, indicate that there is a type node with
//same name already existing
func (node *Script) NewType(expr Expr) (old Expr, ok bool) {
	old, ok = node.Package().NewType(expr)

	if ok {
		node.Types = append(node.Types, expr)
	}

	return
}

//Package override Package function implement
func (node *Script) Package() *Package {
	return node.pkg
}
