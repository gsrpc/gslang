package ast

// Using instruction
type Using struct {
	_Node // mixin _Node
}

// Script .
type Script struct {
	_Node
	Package string            // script's package name
	using   map[string]*Using // using instruction list
}

// NewScript .
func NewScript(name string) *Script {

	script := &Script{
		using: make(map[string]*Using),
	}

	script._init(name)

	return script
}

// Using .
func (script *Script) Using(name string) *Using {
	using := &Using{}

	using._init(name)

	script.using[name] = using

	return using
}
