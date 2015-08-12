package ast

// Module .
type Module struct {
	_Node
	scripts map[string]*Script
	Types   map[string]Type
}

// NewModule create new module
func NewModule(name string) *Module {
	module := &Module{
		scripts: make(map[string]*Script),
	}

	module._init(name)

	return module
}

// Foreach foreach script
func (module *Module) Foreach(f func(script *Script) bool) {
	for _, script := range module.scripts {
		if !f(script) {
			return
		}
	}
}
