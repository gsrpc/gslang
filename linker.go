package gslang

import (
	"fmt"

	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslogger"
)

type _Linker struct {
	gslogger.Log                     //Mixin logger
	types        map[string]ast.Type // defined types
	errorHandler ErrorHandler        // error handler
}

// Link do sematic paring and type link
func (complier *Compiler) Link() {

	linker := &_Linker{
		Log:          gslogger.Get("linker"),
		types:        make(map[string]ast.Type),
		errorHandler: complier.errorHandler,
	}

	for _, script := range complier.scripts {
		linker.createSymbolTable(script)
	}

	for _, script := range complier.scripts {
		linker.createLocalSymbolTable(script)
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

func (linker *_Linker) createLocalSymbolTable(script *ast.Script) {
	linker.D("create script local symoble table : %s", script)
}

func (linker *_Linker) createSymbolTable(script *ast.Script) {

	linker.D("create global symoble table , search script defined types: %s", script)

	script.TableForeach(func(table *ast.Table) {
		fullname := _fullName(script.Package, table)

		if previous, ok := linker.types[fullname]; ok {
			linker.duplicateTypeDef(table, previous)
			return
		}

		linker.types[fullname] = table
	})

	script.ContractForeach(func(contract *ast.Contract) {
		fullname := _fullName(script.Package, contract)

		if previous, ok := linker.types[fullname]; ok {
			linker.duplicateTypeDef(contract, previous)
			return
		}

		linker.types[fullname] = contract
	})

	script.EnumForeach(func(enum *ast.Enum) {
		fullname := _fullName(script.Package, enum)

		if previous, ok := linker.types[fullname]; ok {
			linker.duplicateTypeDef(enum, previous)
			return
		}

		linker.types[fullname] = enum
	})
}
