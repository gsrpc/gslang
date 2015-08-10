package gslang

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslogger"
)

type _Linker struct {
	gslogger.Log                       //Mixin logger
	builtinMapping map[string]string   // builtin types fullname mapping
	types          map[string]ast.Type // defined types
	importTypes    map[string]ast.Type // defined types
	errorHandler   ErrorHandler        // error handler
}

// Link do sematic paring and type link
func (complier *Compiler) Link() {

	linker := &_Linker{
		Log:            gslogger.Get("linker"),
		types:          make(map[string]ast.Type),
		builtinMapping: complier.builtinMapping,
		errorHandler:   complier.errorHandler,
	}

	for _, script := range complier.scripts {
		linker.createSymbolTable(script)
	}

	for _, script := range complier.scripts {
		linker.linkTypes(script)
	}
}

func _fullName(namespace string, node ast.Type) string {
	return fmt.Sprintf("%s.%s", namespace, node)
}

func (linker *_Linker) duplicateTypeDef(lhs ast.Type, rhs ast.Type) {

	start, end := Pos(lhs)

	errinfo := &Error{
		Stage:   StageSemParing,
		Orignal: ErrDuplicateType,
		Start:   start,
		End:     end,
		Text:    fmt.Sprintf("duplicate type defined. see previous type defined here:\n%s", rhs),
	}

	linker.errorHandler.HandleError(errinfo)
}

func (linker *_Linker) unknownTypeRef(node ast.Node, fullname string) {

	start, end := Pos(node)

	errinfo := &Error{
		Stage:   StageSemParing,
		Orignal: ErrTypeNotFound,
		Start:   start,
		End:     end,
		Text:    fmt.Sprintf("not found type reference :%s", fullname),
	}

	linker.errorHandler.HandleError(errinfo)
}

func (linker *_Linker) linkTypes(script *ast.Script) {
	linker.D("create using symbol table for script : %s", script)

	linker.importTypes = make(map[string]ast.Type)

	script.UsingForeach(func(using *ast.Using) {

		name := path.Base(strings.Replace(using.Name(), ".", "/", -1))

		linker.D("link using(%s) : %s", using, name)

		if gslangType, ok := linker.types[using.Name()]; ok {
			linker.importTypes[name] = gslangType
			linker.D("link using(%s) : %s -- success", using, name)
			return
		}

		linker.unknownTypeRef(using, using.Name())
	})

	linker.D("create using symbol table for script : %s -- success", script)

	script.TypeForeach(func(gslangType ast.Type) {

		linker.linkType(script, gslangType)
	})
}

func (linker *_Linker) linkType(script *ast.Script, gslangType ast.Type) {
	fullname := _fullName(script.Package, gslangType)

	linker.D("link type(%s) :%s", fullname, reflect.TypeOf(gslangType))

	switch gslangType.(type) {
	case *ast.Table:
		linker.linkTable(script, gslangType.(*ast.Table))
	case *ast.Contract:
		linker.linkContract(script, gslangType.(*ast.Contract))
	case *ast.Enum:
	case *ast.TypeRef:
		linker.linkTypeRef(script, gslangType.(*ast.TypeRef))
	case *ast.Seq:
		linker.linkType(script, gslangType.(*ast.Seq).Component)
	}
}

func (linker *_Linker) linkTypeRef(script *ast.Script, typeRef *ast.TypeRef) {

	linker.D("try link type reference :%s", typeRef)

	linkedType, ok := script.Type(typeRef.Name())

	if ok {
		typeRef.Ref = linkedType

		linker.D("found type %s", linkedType)

		return
	}

	linkedType, ok = linker.importTypes[typeRef.Name()]

	if ok {
		typeRef.Ref = linkedType

		linker.D("found import types %s", linkedType)

		return
	}

	linker.unknownTypeRef(typeRef, typeRef.Name())
}

func (linker *_Linker) linkContract(script *ast.Script, contract *ast.Contract) {

	linker.D("link contract(%s)", contract)

	for _, method := range contract.Methods {
		linker.D("link method(%s) return type(%s)", method, method.Return)
		linker.linkType(script, method.Return)
		linker.D("link method(%s) return type(%s) -- success", method, method.Return)

		for _, param := range method.Params {
			linker.D("link method(%s) param type(%s)", method, param)
			linker.linkType(script, param.Type)
			linker.D("link method(%s) param type(%s) -- success", method, param.Type)
		}
	}

	linker.D("link contract(%s) -- success", contract)
}

func (linker *_Linker) linkTable(script *ast.Script, table *ast.Table) {
	for _, field := range table.Fields {
		linker.linkType(script, field.Type)
	}
}

func (linker *_Linker) createSymbolTable(script *ast.Script) {

	linker.D("create global symoble table , search script defined types: %s", script)

	script.TypeForeach(func(gslangType ast.Type) {
		fullname := _fullName(script.Package, gslangType)

		if previous, ok := linker.types[fullname]; ok {
			linker.duplicateTypeDef(gslangType, previous)
			return
		}

		linker.types[fullname] = gslangType
	})
}
