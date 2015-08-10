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
		}

		linker.unknownTypeRef(using, using.Name())
	})

	linker.D("create using symbol table for script : %s -- success", script)

	script.TypeForeach(func(gslangType ast.Type) {

		fullname := _fullName(script.Package, gslangType)

		linker.D("link type(%s) :%s", fullname, reflect.TypeOf(gslangType))

		switch gslangType.(type) {
		case *ast.Table:
			linker.linkTable(script, gslangType.(*ast.Table))
		case *ast.Contract:

		case *ast.Enum:
		}

		// for _, field := range table.Fields {
		// 	linker.D("link field(%s) : %s", field, reflect.TypeOf(field.Type))
		//
		// 	if typeRef, ok := field.Type.(*ast.TypeRef); ok {
		// 		typeRef.Name()
		// 	}
		// }
	})
}

func (linker *_Linker) linkTable(script *ast.Script, table *ast.Table) {
	for _, field := range table.Fields {
		if typeRef, ok := field.Type.(*ast.TypeRef); ok {
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
		}
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
