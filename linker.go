package gslang

import (
	"fmt"
	"path"
	"strings"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslogger"
)

type _Linker struct {
	gslogger.Log                     //Mixin logger
	types        map[string]ast.Type // defined types
	importTypes  map[string]ast.Type // defined types
	errorHandler ErrorHandler        // error handler
	linkdepth    int                 // link depth
	compiler     *Compiler           // compiler
}

// Link do sematic paring and type link
func (compiler *Compiler) Link() (err error) {

	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	linker := &_Linker{
		Log:          gslogger.Get("linker"),
		types:        make(map[string]ast.Type),
		errorHandler: compiler.errorHandler,
		compiler:     compiler,
	}

	compiler.module.Foreach(func(script *ast.Script) bool {
		linker.createSymbolTable(script)
		return true
	})

	compiler.module.Types = linker.types

	compiler.module.Foreach(func(script *ast.Script) bool {
		linker.linkTypes(script)
		return true
	})

	compiler.module.Foreach(func(script *ast.Script) bool {

		script.TypeForeach(func(gslangType ast.Type) {

			linker.checkAnnotation(script, gslangType)
		})

		return true
	})

	return
}

func _fullName(namespace string, node ast.Type) string {
	return fmt.Sprintf("%s.%s", namespace, node)
}

func (linker *_Linker) Eval() Eval {
	return linker.compiler.Eval()
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

func (linker *_Linker) errorf(err error, node ast.Node, fmtstr string, args ...interface{}) {
	start, end := Pos(node)

	errinfo := &Error{
		Stage:   StageSemParing,
		Orignal: err,
		Start:   start,
		End:     end,
		Text:    fmt.Sprintf(fmtstr, args...),
	}

	linker.errorHandler.HandleError(errinfo)
}

func (linker *_Linker) startLinkNode(node ast.Node) {

	linker.D("<link %s>", node)

	linker.linkdepth++
}

func (linker *_Linker) endLinkNode(node ast.Node) {
	linker.linkdepth--

	linker.D("</link %s>", node)
}

func (linker *_Linker) D(fmtstr string, args ...interface{}) {
	linker.Log.D("%s%s", strings.Repeat(" ", linker.linkdepth*2), fmt.Sprintf(fmtstr, args...))
}

func (linker *_Linker) linkTypes(script *ast.Script) {
	linker.D("create using symbol table for script : %s", script)

	linker.importTypes = make(map[string]ast.Type)

	script.UsingForeach(func(using *ast.Using) {

		name := path.Base(strings.Replace(using.Name(), ".", "/", -1))

		linker.D("link using(%s) : %s", using, name)

		if gslangType, ok := linker.types[using.Name()]; ok {
			linker.importTypes[name] = gslangType
			using.Ref = gslangType
			linker.D("link using(%s:%p) : %s -- success", using, using, name)
			return
		}

		linker.errorf(ErrTypeNotFound, using, "using statment reference unknown type(%s)", using.Name())
	})

	linker.D("create using symbol table for script : %s -- success", script)

	script.TypeForeach(func(gslangType ast.Type) {

		linker.linkType(script, gslangType)
	})
}

func (linker *_Linker) linkTypeAnnotation(script *ast.Script, typeDecl ast.Type) {
	for _, annotation := range Annotations(typeDecl) {
		linker.linkAnnotation(script, annotation)
	}
}

func (linker *_Linker) checkAnnotation(script *ast.Script, typeDecl ast.Type) {
	for _, annotation := range Annotations(typeDecl) {

		if annotation.Type.Ref == nil {
			continue
		}

		usage, ok := FindAnnotation(annotation.Type.Ref, "gslang.annotations.Usage")

		if !ok {
			linker.errorf(ErrAnnotation, annotation, "illegal annotation type : table(%s) must be annotation by gslang.annotations.Usage", annotation.Type.Ref.FullName())
		}

		if usage.Args.Count() != 1 {
			continue
		}

		val := linker.Eval().EvalInt(usage.Args.Arg(0))

		scriptConstant := int64(linker.Eval().EvalEnumConstant("gslang.annotations.Target", "Script"))

		moduleConstant := int64(linker.Eval().EvalEnumConstant("gslang.annotations.Target", "Module"))

		if scriptConstant&val != 0 {
			linker.D("move anntotation(%s) to script(%s)", annotation, script)
			_RemoveAnnotation(typeDecl, annotation)
			_AttachAnnotation(script, annotation)
		} else if moduleConstant&val != 0 {
			linker.D("move anntotation(%s) to module(%s)", annotation, script.Module)
			_RemoveAnnotation(typeDecl, annotation)
			_AttachAnnotation(script.Module, annotation)
		}
	}

	switch typeDecl.(type) {
	case *ast.Contract:
		linker.checkContractAnnotation(script, typeDecl.(*ast.Contract))
	}
}

func (linker *_Linker) checkContractAnnotation(script *ast.Script, contract *ast.Contract) {

	for _, method := range contract.Methods {
		for _, exception := range method.Exceptions {
			ref, ok := exception.Type.(*ast.TypeRef)

			if !ok {
				linker.errorf(ErrType, exception, "exception must be table with @Exception annotation")
				continue
			}

			if ref.Ref != nil {
				if _, ok := ref.Ref.(*ast.Table); !ok {
					linker.errorf(ErrType, exception, "exception must be table with @Exception annotation")
					continue
				}

				_, ok = FindAnnotation(ref.Ref, "gslang.Exception")

				if !ok {
					linker.errorf(ErrType, exception, "exception must be table with @Exception annotation")
				}
			}
		}
	}

}

func (linker *_Linker) linkType(script *ast.Script, gslangType ast.Type) {

	linker.startLinkNode(gslangType)

	switch gslangType.(type) {
	case *ast.Table:
		linker.linkTypeAnnotation(script, gslangType)
		linker.linkTable(script, gslangType.(*ast.Table))
	case *ast.Contract:
		linker.linkTypeAnnotation(script, gslangType)
		linker.linkContract(script, gslangType.(*ast.Contract))
	case *ast.Enum:
		linker.linkTypeAnnotation(script, gslangType)
		linker.linkEnum(script, gslangType.(*ast.Enum))
	case *ast.TypeRef:
		linker.linkTypeRef(script, gslangType.(*ast.TypeRef))
	case *ast.Seq:
		linker.linkType(script, gslangType.(*ast.Seq).Component)
	}

	linker.endLinkNode(gslangType)
}

func (linker *_Linker) linkTypeRef(script *ast.Script, typeRef *ast.TypeRef) {

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

	linker.errorf(ErrTypeNotFound, typeRef, "unknown type reference :%s", typeRef)
}

func (linker *_Linker) linkExpr(script *ast.Script, expr ast.Expr) {

	gserrors.Assert(expr != nil, "input arg expr can't be nil")

	linker.startLinkNode(expr)

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

	linker.endLinkNode(expr)
}

func (linker *_Linker) linkUnaryOp(script *ast.Script, unary *ast.UnaryOp) {

	linker.linkExpr(script, unary.Operand)
}

func (linker *_Linker) linkBinaryOp(script *ast.Script, binary *ast.BinaryOp) {

	linker.linkExpr(script, binary.LHS)

	linker.linkExpr(script, binary.RHS)
}

func (linker *_Linker) linkNewObj(script *ast.Script, newObj *ast.NewObj) {
	linker.linkType(script, newObj.Type)

	if newObj.Args != nil {
		linker.linkExpr(script, newObj.Args)

		if newObj.Type.Ref != nil {
			switch newObj.Type.Ref.(type) {
			case *ast.Table:
				linker.linkTableNewObj(script, newObj.Type.Ref.(*ast.Table), newObj.Args)
			}
		}
	}
}

func (linker *_Linker) linkTableNewObj(script *ast.Script, table *ast.Table, args *ast.ArgsTable) {
	if args.Named {

		for _, arg := range args.Args() {

			namedArg := arg.(*ast.NamedArg)

			_, ok := table.Field(namedArg.Name())

			if !ok {
				linker.errorf(ErrFieldName, arg, "unknown table(%s) field(%s)", table, namedArg)
			}

			//TODO : check if field type match the arg expr
		}

		return
	}

	if len(table.Fields) != args.Count() {
		linker.errorf(ErrNewObj, args, "wrong newobj args num for table(%s) : expect %d but got %d", table, len(table.Fields), args.Count())

		//TODO : check if field type match the arg expr
	}
}

func (linker *_Linker) linkArgsTable(script *ast.Script, argsTable *ast.ArgsTable) {
	for _, arg := range argsTable.Args() {
		linker.linkExpr(script, arg)
	}
}

func (linker *_Linker) linkNameArg(script *ast.Script, namedArg *ast.NamedArg) {
	linker.linkExpr(script, namedArg.Arg)
}

func (linker *_Linker) linkConstantRef(script *ast.Script, constantRef *ast.ConstantRef) {

	nodes := strings.Split(constantRef.Name(), ".")

	if len(nodes) < 2 {
		linker.errorf(ErrTypeNotFound, constantRef, "unknown constant val (%s)", constantRef.Name())
		return
	}

	name := nodes[len(nodes)-1]

	typeRef := ast.NewTypeRef(strings.Join(nodes[0:len(nodes)-1], "."))

	start, end := Pos(constantRef)

	_setNodePos(typeRef, start, end)

	linker.linkTypeRef(script, typeRef)

	if typeRef.Ref != nil {
		enum := typeRef.Ref.(*ast.Enum)

		for _, constant := range enum.Constants {

			if constant.Name() == name {
				constantRef.Value = constant
				return
			}
		}

		linker.errorf(ErrVariableName, constantRef, "unknown enum(%s) constant filed :%s", enum, name)
	}

}

func (linker *_Linker) linkAnnotation(script *ast.Script, annotation *ast.Annotation) {
	linker.linkNewObj(script, ast.NewNewObj2(annotation.Type, annotation.Args))
}

func (linker *_Linker) linkEnum(script *ast.Script, enum *ast.Enum) {

}

func (linker *_Linker) linkContract(script *ast.Script, contract *ast.Contract) {

	for _, method := range contract.Methods {
		for _, annotation := range Annotations(method) {
			linker.linkAnnotation(script, annotation)
		}

		linker.linkType(script, method.Return)

		// link method params
		for _, param := range method.Params {

			for _, annotation := range Annotations(param) {

				linker.linkAnnotation(script, annotation)
			}

			linker.linkType(script, param.Type)
		}

		for _, exception := range method.Exceptions {
			for _, annotation := range Annotations(exception) {

				linker.linkAnnotation(script, annotation)
			}

			linker.linkType(script, exception.Type)
		}
	}
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

	linker.D("create global symoble table , search script defined types: %s -- success", script)
}
