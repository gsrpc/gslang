package ast

import (
	"fmt"

	"github.com/gsrpc/gslang/lexer"
)

// Type .
type Type interface {
	Node
	FullName() string
	Package() string
	Script() string
}

// TypeDecl .
type TypeDecl interface {
	Type
	Module() *Module
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

// FullName .
func (ref *TypeRef) FullName() string {
	if ref.Ref != nil {
		return "ref :" + ref.Ref.FullName()
	}

	return "unlink ref :" + ref.Name()
}

// Package .
func (ref *TypeRef) Package() string {
	return ""
}

// Script .
func (ref *TypeRef) Script() string {
	return ""
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

// FullName .
func (builtin *BuiltinType) FullName() string {
	return builtin.Name()
}

// Package .
func (builtin *BuiltinType) Package() string {
	return ""
}

// Script .
func (builtin *BuiltinType) Script() string {
	return ""
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
	script *Script  // script belongs to
}

// NewTable .
func (script *Script) NewTable(name string) (Type, bool) {

	if table, ok := script.types[name]; ok {
		return table, false
	}

	table := &Table{
		script: script,
	}

	table._init(name)

	script.types[name] = table

	return table, true
}

// Module .
func (table *Table) Module() *Module {
	return table.script.Module
}

// FullName .
func (table *Table) FullName() string {
	return table.script.Package + "." + table.Name()
}

// Package .
func (table *Table) Package() string {
	return table.script.Package
}

// Script .
func (table *Table) Script() string {
	return table.script.String()
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
	ID   int
	Type Type
}

// Exception .
type Exception struct {
	_Node
	Type Type
	ID   int8
}

// Method .
type Method struct {
	_Node
	ID         int          /// id
	Return     Type         // return type
	Params     []*Param     // Params type list
	Exceptions []*Exception // exception list
}

// ParamsCount .
func (method *Method) ParamsCount() int {
	return len(method.Params)
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
func (method *Method) NewException(typeDecl Type) *Exception {

	exception := &Exception{
		Type: typeDecl,
		ID:   int8(len(method.Exceptions)),
	}

	exception._init(typeDecl.Name())

	method.Exceptions = append(method.Exceptions, exception)

	return exception
}

// NewParam .
func (method *Method) NewParam(name string, typeDecl Type) (*Param, bool) {
	if param, ok := method.Param(name); ok {
		return param, false
	}

	param := &Param{
		ID:   len(method.Params),
		Type: typeDecl,
	}

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
	script    *Script
}

// NewEnum .
func (script *Script) NewEnum(name string) (Type, bool) {
	if enum, ok := script.types[name]; ok {
		return enum, false
	}

	enum := &Enum{
		script: script,
	}

	enum._init(name)

	script.types[name] = enum

	return enum, true
}

// Package .
func (enum *Enum) Package() string {
	return enum.script.Package
}

// Module .
func (enum *Enum) Module() *Module {
	return enum.script.Module
}

// FullName .
func (enum *Enum) FullName() string {
	return enum.script.Package + "." + enum.Name()
}

// Script .
func (enum *Enum) Script() string {
	return enum.script.String()
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
	} else {
		constant.Value = enum.Constants[len(enum.Constants)-1].Value + 1
	}

	enum.Constants = append(enum.Constants, constant)

	return constant, true
}

// Contract .
type Contract struct {
	_Node             // Mixin default node implement
	Methods []*Method // table fields
	script  *Script
}

// NewContract .
func (script *Script) NewContract(name string) (Type, bool) {

	if contract, ok := script.types[name]; ok {
		return contract, false
	}

	contract := &Contract{
		script: script,
	}

	contract._init(name)

	script.types[name] = contract

	return contract, true
}

// Package .
func (contract *Contract) Package() string {
	return contract.script.Package
}

// Module .
func (contract *Contract) Module() *Module {
	return contract.script.Module
}

// FullName .
func (contract *Contract) FullName() string {
	return contract.script.Package + "." + contract.Name()
}

// Script .
func (contract *Contract) Script() string {
	return contract.script.String()
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

	method := &Method{
		ID: len(contract.Methods),
	}

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

	seq._init(fmt.Sprintf("%s[%d]", component, size))

	return seq
}

// FullName .
func (seq *Seq) FullName() string {

	if seq.Size > 0 {
		return fmt.Sprintf("%s[%d]", seq.Component.FullName(), seq.Size)
	}

	return fmt.Sprintf("%s[]", seq.Component.FullName())
}

// Package .
func (seq *Seq) Package() string {
	return "gslang"
}

// Script .
func (seq *Seq) Script() string {
	return "gslang.gs"
}
