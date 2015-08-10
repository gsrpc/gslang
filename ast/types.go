package ast

import (
	"fmt"

	"github.com/gsdocker/gslang/lexer"
)

// Type .
type Type interface {
	Node
}

// TypeRef .
type TypeRef struct {
	_Node
	Ref Type
}

// NewTypeRef .
func NewTypeRef(name string) *TypeRef {
	typeRef := &TypeRef{}

	typeRef._init(name)

	return typeRef
}

// BuiltinType .
type BuiltinType struct {
	_Node
	Type lexer.TokenType
}

// NewBuiltinType .
func NewBuiltinType(builtin lexer.TokenType) *BuiltinType {
	builtinType := &BuiltinType{
		Type: builtin,
	}

	builtinType._init(builtin.String())

	return builtinType
}

// Field .
type Field struct {
	_Node
	Type Type // Field Type
}

// Table .
type Table struct {
	_Node           // Mixin default node implement
	Fields []*Field // table fields
	Script *Script  // script belongs to
}

// NewTable .
func (script *Script) NewTable(name string) (Type, bool) {

	if table, ok := script.types[name]; ok {
		return table, false
	}

	table := &Table{
		Script: script,
	}

	table._init(name)

	script.types[name] = table

	return table, true
}

// Field .
func (table *Table) Field(name string) (*Field, bool) {
	for _, field := range table.Fields {
		if field.Name() == name {
			return field, true
		}
	}

	return nil, false
}

// NewField .
func (table *Table) NewField(name string, typeDecl Type) (*Field, bool) {
	if field, ok := table.Field(name); ok {
		return field, false
	}

	field := &Field{Type: typeDecl}

	field._init(name)

	table.Fields = append(table.Fields, field)

	return field, true
}

// Param .
type Param struct {
	_Node
	Type Type
}

// Method .
type Method struct {
	_Node
	Return     Type     // return type
	Params     []*Param // Params type list
	Exceptions []Type   // exception list
}

// Param .
func (method *Method) Param(name string) (*Param, bool) {

	for _, param := range method.Params {
		if param.Name() == name {
			return param, true
		}
	}

	return nil, false
}

// NewException .
func (method *Method) NewException(typeDecl Type) {
	method.Exceptions = append(method.Exceptions, typeDecl)
}

// NewParam .
func (method *Method) NewParam(name string, typeDecl Type) (*Param, bool) {
	if param, ok := method.Param(name); ok {
		return param, false
	}

	param := &Param{Type: typeDecl}

	param._init(name)

	method.Params = append(method.Params, param)

	return param, true
}

// EnumConstant .
type EnumConstant struct {
	_Node       // Mixin default node implement
	Value int32 // constant value
}

// Enum .
type Enum struct {
	_Node                     // Mixin default node implement
	Constants []*EnumConstant // table fields
}

// NewEnum .
func (script *Script) NewEnum(name string) (Type, bool) {
	if enum, ok := script.types[name]; ok {
		return enum, false
	}

	enum := &Enum{}

	enum._init(name)

	script.types[name] = enum

	return enum, true
}

// Constant .
func (enum *Enum) Constant(name string) (*EnumConstant, bool) {
	for _, constant := range enum.Constants {
		if constant.Name() == name {
			return constant, true
		}
	}

	return nil, false
}

// NewConstant .
func (enum *Enum) NewConstant(name string) (*EnumConstant, bool) {
	if constant, ok := enum.Constant(name); ok {
		return constant, false
	}

	constant := &EnumConstant{}

	constant._init(name)

	if len(enum.Constants) == 0 {
		constant.Value = 0
		return constant, true
	}

	constant.Value = enum.Constants[len(enum.Constants)-1].Value + 1

	return constant, true
}

// Contract .
type Contract struct {
	_Node             // Mixin default node implement
	Methods []*Method // table fields
}

// NewContract .
func (script *Script) NewContract(name string) (Type, bool) {

	if contract, ok := script.types[name]; ok {
		return contract, false
	}

	contract := &Contract{}

	contract._init(name)

	script.types[name] = contract

	return contract, true
}

// Method .
func (contract *Contract) Method(name string) (*Method, bool) {
	for _, method := range contract.Methods {
		if method.Name() == name {
			return method, true
		}
	}

	return nil, false
}

// NewMethod .
func (contract *Contract) NewMethod(name string) (*Method, bool) {
	if method, ok := contract.Method(name); ok {
		return method, false
	}

	method := &Method{}

	method._init(name)

	contract.Methods = append(contract.Methods, method)

	return method, true
}

// Seq Type seq
type Seq struct {
	_Node
	Component Type
	Size      int
}

// NewSeq .
func NewSeq(component Type, size int) *Seq {
	seq := &Seq{
		Component: component,
		Size:      size,
	}

	seq._init(fmt.Sprintf("seq<%s>(%d)", component, size))

	return seq
}
