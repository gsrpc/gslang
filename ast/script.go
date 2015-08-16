package ast

// Using instruction
type Using struct {
	_Node      // mixin _Node
	Ref   Type // type reference
}

// Script .
type Script struct {
	_Node
	Package string            // script's package name
	using   map[string]*Using // using instruction list
	types   map[string]Type   // tables
}

// NewScript .
func (module *Module) NewScript(name string) *Script {

	script := &Script{
		using: make(map[string]*Using),
		types: make(map[string]Type),
	}

	script._init(name)

	module.scripts[name] = script

	return script
}

// UsingForeach .
func (script *Script) UsingForeach(f func(*Using)) {
	for _, using := range script.using {
		f(using)
	}
}

// TypeForeach .
func (script *Script) TypeForeach(f func(Type)) {
	for _, gslangType := range script.types {
		f(gslangType)
	}
}

// Using .
func (script *Script) Using(name string) *Using {
	using := &Using{}

	using._init(name)

	script.using[name] = using

	return using
}

// Type get type .
func (script *Script) Type(name string) (Type, bool) {
	gslangType, ok := script.types[name]

	return gslangType, ok
}
