package gslang

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"github.com/gsdocker/gserrors"
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
func (complier *Compiler) Link() (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

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

	return
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
		Text:    fmt.Sprintf("unknown type reference :%s", fullname),
	}

	linker.errorHandler.HandleError(errinfo)
}

func (linker *_Linker) unknownConstant(node ast.Node, constantType ast.Type, name string) {

	start, end := Pos(node)

	errinfo := &Error{
		Stage:   StageSemParing,
		Orignal: ErrTypeNotFound,
		Start:   start,
		End:     end,
		Text:    fmt.Sprintf("unknown constant(%s.%s)", constantType, name),
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

	linker.D("link type(%s) :%s", gslangType, reflect.TypeOf(gslangType))

	switch gslangType.(type) {
	case *ast.Table:
		linker.linkTable(script, gslangType.(*ast.Table))
	case *ast.Contract:
		linker.linkContract(script, gslangType.(*ast.Contract))
	case *ast.Enum:
		linker.linkEnum(script, gslangType.(*ast.Enum))
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

func (linker *_Linker) linkExpr(script *ast.Script, expr ast.Expr) {

	gserrors.Assert(expr != nil, "input arg expr can't be nil")

	linker.D("link expr(%s) :%s", expr, reflect.TypeOf(expr))

	switch expr.(type) {
	case *ast.ArgsTable:
		linker.linkArgsTable(script, expr.(*ast.ArgsTable))
	case *ast.NamedArg:
		linker.linkNameArg(script, expr.(*ast.NamedArg))
	case *ast.ConstantRef:
		linker.linkConstantRef(script, expr.(*ast.ConstantRef))
	case *ast.NewObj:
		linker.linkNewObj(script, expr.(*ast.NewObj))
	case *ast.UnaryOp:
		linker.linkUnaryOp(script, expr.(*ast.UnaryOp))
	case *ast.BinaryOp:
		linker.linkBinaryOp(script, expr.(*ast.BinaryOp))
	}

	linker.D("link expr(%s) :%s -- success", expr, reflect.TypeOf(expr))
}

func (linker *_Linker) linkUnaryOp(script *ast.Script, unary *ast.UnaryOp) {
	linker.D("link binary op %s%s", unary, unary.Operand)

	linker.linkExpr(script, unary.Operand)

	linker.D("link binary op %s%s", unary, unary.Operand)
}

func (linker *_Linker) linkBinaryOp(script *ast.Script, binary *ast.BinaryOp) {

	linker.D("link binary op %s%s%s", binary.LHS, binary, binary.RHS)

	linker.linkExpr(script, binary.LHS)

	linker.linkExpr(script, binary.RHS)

	linker.D("link binary op %s%s%s -- success", binary.LHS, binary, binary.RHS)
}

func (linker *_Linker) linkNewObj(script *ast.Script, newObj *ast.NewObj) {
	linker.D("link newobj(%s)", newObj)
	linker.linkType(script, newObj.Type)

	if newObj.Args != nil {
		linker.linkExpr(script, newObj.Args)
	}

	linker.D("link newobj(%s) -- success", newObj)
}

func (linker *_Linker) linkArgsTable(script *ast.Script, argsTable *ast.ArgsTable) {
	for _, arg := range argsTable.GetArgs() {
		linker.linkExpr(script, arg)
	}
}

func (linker *_Linker) linkNameArg(script *ast.Script, namedArg *ast.NamedArg) {
	linker.D("link name arg(%s)", namedArg)
	linker.linkExpr(script, namedArg.Arg)
	linker.D("link name arg(%s) -- success", namedArg)
}

func (linker *_Linker) linkConstantRef(script *ast.Script, constantRef *ast.ConstantRef) {
	linker.D("link constant reference(%s)", constantRef)

	nodes := strings.Split(constantRef.Name(), ".")

	if len(nodes) == 0 {
		linker.unknownTypeRef(constantRef, constantRef.Name())
		return
	}

	name := nodes[len(nodes)-1]

	typeRef := ast.NewTypeRef(strings.Join(nodes[0:len(nodes)-1], "."))

	linker.linkTypeRef(script, typeRef)

	if typeRef.Ref != nil {
		enum := typeRef.Ref.(*ast.Enum)

		for _, constant := range enum.Constants {

			if constant.Name() == name {
				linker.D("link constant reference(%s) -- success", constantRef)
				return
			}
		}

		linker.unknownConstant(constantRef, typeRef.Ref, name)
	}

}

func (linker *_Linker) linkAnnotation(script *ast.Script, annotation *ast.Annotation) {

	linker.D("link annotation(%s)", annotation)

	linker.linkNewObj(script, ast.NewNewObj(annotation.Type.Name(), annotation.Args))

	linker.D("link annotation(%s) -- success", annotation)
}

func (linker *_Linker) linkEnum(script *ast.Script, enum *ast.Enum) {
	linker.D("link enum(%s)", enum)

	for _, annotation := range Annotation(enum) {
		linker.linkAnnotation(script, annotation)
	}

	linker.D("link enum(%s) -- success", enum)
}

func (linker *_Linker) linkContract(script *ast.Script, contract *ast.Contract) {

	linker.D("link contract(%s)", contract)

	for _, annotation := range Annotation(contract) {

		linker.D("link contract(%s) annotation(%s)", contract, annotation)

		linker.linkAnnotation(script, annotation)

		linker.D("link contract(%s) annotation(%s) -- sucess", contract, annotation)
	}

	for _, method := range contract.Methods {
		for _, annotation := range Annotation(method) {

			linker.D("link method(%s) annotation(%s)", method, annotation)

			linker.linkAnnotation(script, annotation)

			linker.D("link method(%s) annotation(%s) -- sucess", method, annotation)
		}

		linker.D("link method(%s) return type(%s)", method, method.Return)
		linker.linkType(script, method.Return)
		linker.D("link method(%s) return type(%s) -- success", method, method.Return)

		for _, param := range method.Params {

			for _, annotation := range Annotation(param) {

				linker.D("link method(%s) param(%s) annotation(%s)", method, param, annotation)

				linker.linkAnnotation(script, annotation)

				linker.D("link method(%s) param(%s) annotation(%s) -- sucess", method, param, annotation)
			}

			linker.D("link method(%s) param type(%s)", method, param)
			linker.linkType(script, param.Type)
			linker.D("link method(%s) param type(%s) -- success", method, param.Type)
		}
	}

	linker.D("link contract(%s) -- success", contract)
}

func (linker *_Linker) linkTable(script *ast.Script, table *ast.Table) {

	for _, annotation := range Annotation(table) {

		linker.D("link table(%s) annotation(%s)", table, annotation)

		linker.linkAnnotation(script, annotation)

		linker.D("link table(%s) annotation(%s) -- sucess", table, annotation)
	}

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
