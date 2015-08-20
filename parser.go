package gslang

import (
	"bytes"
	"fmt"

	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslang/lexer"
	"github.com/gsdocker/gslogger"
)

// Parser gslang parser
type Parser struct {
	gslogger.Log                      // mixin logger
	lexer           *lexer.Lexer      // token lexer
	script          *ast.Script       // current script
	commentStack    []*ast.Comment    // comment stack
	annotationStack []*ast.Annotation // comment stack
	errorHandler    ErrorHandler      // error Handler
}

func (compiler *Compiler) parse(lexer *lexer.Lexer, errorHandler ErrorHandler) *ast.Script {

	return (&Parser{
		Log:          gslogger.Get("parser"),
		lexer:        lexer,
		script:       compiler.module.NewScript(lexer.String()),
		errorHandler: errorHandler,
	}).parse()
}

func (parser *Parser) peek() *lexer.Token {
	token, err := parser.lexer.Peek()
	if err != nil {
		panic(err)
	}

	return token
}

func (parser *Parser) next() (token *lexer.Token) {
	token, err := parser.lexer.Next()
	if err != nil {
		panic(err)
	}

	return token
}

func (parser *Parser) errorf(position lexer.Position, fmtstring string, args ...interface{}) {

	errinfo := &Error{
		Stage:   StageParing,
		Orignal: ErrParser,
		Start:   position,
		End:     position,
		Text:    fmt.Sprintf(fmtstring, args...),
	}

	parser.errorHandler.HandleError(errinfo)
}

func (parser *Parser) errorf2(err error, position lexer.Position, fmtstring string, args ...interface{}) {

	errinfo := &Error{
		Stage:   StageParing,
		Orignal: err,
		Start:   position,
		End:     position,
		Text:    fmt.Sprintf(fmtstring, args...),
	}

	parser.errorHandler.HandleError(errinfo)
}

func (parser *Parser) expectf(expect lexer.TokenType, fmtstring string, args ...interface{}) *lexer.Token {

	for {
		token := parser.next()

		if token.Type != expect {
			parser.errorf(token.Start, fmt.Sprintf("current token(%s) \n%s", token.Type, fmt.Sprintf(fmtstring, args...)))
			continue
		}

		return token
	}
}

func (parser *Parser) parse() *ast.Script {

	parser.parsePackage()

	// parse import instructions

	for parser.parseImport() {
	}

	for parser.parseType() {

	}

	parser.attachComment(parser.script)

	parser.attachAnnotation(parser.script)

	return parser.script
}

func (parser *Parser) parseType() bool {
	for parser.parseAnnotation() {

	}

	token := parser.peek()

	switch token.Type {
	case lexer.KeyTable:
		parser.expectTable("expect table type define")
		return true
	case lexer.KeyContract:
		parser.expectContract("expect contract type define")
		return true
	case lexer.KeyEnum:
		parser.expectEnum("expect enum type define")
		return true
	case lexer.TokenEOF:
		return false
	default:
		parser.errorf(token.Start, "unexpect token\n%s", token)
	}
	return false
}

func (parser *Parser) parsePackage() {

	parser.D("[] parse script's package line")

	for parser.parseComment() {
	}

	parser.expectf(lexer.KeyPackage, "script must start with package keyword")

	parser.script.Package, _, _ = parser.expectFullName("expect script's package name")

	parser.expectf(lexer.TokenType(';'), "package name must end with ';'")

	parser.D("package [%s]", parser.script.Package)

	parser.D("parse script's package line -- success")
}

func (parser *Parser) expectEnum(fmtstring string, args ...interface{}) *ast.Enum {

	msg := fmt.Sprintf(fmtstring, args...)

	start := parser.expectf(lexer.KeyEnum, "expect keyword enum").Start

	token := parser.expectf(lexer.TokenID, "expect contract name")

	name := token.Value.(string)

	enum, ok := parser.script.NewEnum(name)

	parser.D("parse enum %s", name)

	parser.attachAnnotation(enum)

	if !ok {
		parser.errorf(token.Start, "%s\n\tduplicate enum(%s) defined", msg, name)
	}

	parser.expectf(lexer.TokenType('{'), "contract body must start with {")

	for {

		for parser.parseComment() {

		}

		token = parser.peek()

		if token.Type == lexer.TokenType('}') {
			break
		}

		constantName := parser.expectf(lexer.TokenID, "expect enum constant name")

		name = constantName.Value.(string)

		constant, ok := enum.(*ast.Enum).NewConstant(name)

		if !ok {
			parser.errorf(token.Start, "%s\n\tduplicate enum(%s) contract(%s) defined", msg, enum, name)
		}

		end := constantName.End

		token = parser.peek()

		if token.Type == lexer.TokenType('(') {

			parser.next()

			val := int32(parser.expectf(lexer.TokenINT, "expect constant value").Value.(int64))

			constant.Value = val

			end = parser.expectf(lexer.TokenType(')'), "enum constant val must end with )").End
		}

		_setNodePos(constant, constantName.Start, end)

		parser.attachComment(constant)

		token = parser.peek()

		if token.Type != lexer.TokenType(',') {
			break
		}

		parser.next()
	}

	end := parser.expectf(lexer.TokenType('}'), "contract body must end with }").End

	_setNodePos(enum, start, end)

	parser.attachComment(enum)

	return enum.(*ast.Enum)

}

func (parser *Parser) expectContract(fmtstring string, args ...interface{}) *ast.Contract {
	msg := fmt.Sprintf(fmtstring, args...)

	start := parser.expectf(lexer.KeyContract, "expect keyword contract").Start

	token := parser.expectf(lexer.TokenID, "expect contract name")

	name := token.Value.(string)

	parser.expectf(lexer.TokenType('{'), "contract body must start with {")

	contract, ok := parser.script.NewContract(name)

	parser.D("parse contract %s", name)

	parser.attachAnnotation(contract)

	if !ok {
		parser.errorf(token.Start, "%s\n\tduplicate contract(%s) defined", msg, name)
	}

	for parser.parseMethodDecl(contract.(*ast.Contract)) {
	}

	end := parser.expectf(lexer.TokenType('}'), "contract body must end with }").End

	_setNodePos(contract, start, end)

	parser.attachComment(contract)

	parser.D("parse contract %s -- success", name)

	return contract.(*ast.Contract)
}

func (parser *Parser) expectTable(fmtstring string, args ...interface{}) *ast.Table {

	msg := fmt.Sprintf(fmtstring, args...)

	start := parser.expectf(lexer.KeyTable, "expect keyword table").Start

	token := parser.expectf(lexer.TokenID, "expect table name")

	name := token.Value.(string)

	parser.expectf(lexer.TokenType('{'), "table body must start with {")

	table, ok := parser.script.NewTable(name)

	parser.attachAnnotation(table)

	parser.D("parse table %s", name)

	if !ok {
		parser.errorf(token.Start, "%s\n\tduplicate table(%s) defined", msg, name)
	}

	for parser.parseFieldDecl(table.(*ast.Table)) {

	}

	end := parser.expectf(lexer.TokenType('}'), "table body must end with }").End

	_setNodePos(table, start, end)

	parser.attachComment(table)

	parser.D("parse table %s -- success", name)

	return table.(*ast.Table)
}

func (parser *Parser) attachAnnotation(node ast.Node) {

	if parser.annotationStack != nil {
		_AttachAnnotation(node, parser.annotationStack)
		parser.D("attach annotations %v to %s", parser.annotationStack, node)
		parser.annotationStack = nil
	}
}

func (parser *Parser) attachComment(node ast.Node) {

	if comment, ok := parser.tailComment(); ok {
		if _AttachComment(node, comment) {
			parser.D("attach %s comment :\n|%s|", node, comment)
			parser.popTailComment()
		}
	}

	for parser.parseComment() {
		comment, _ := parser.tailComment()
		if _AttachComment(node, comment) {
			parser.D("attach %s comment :\n|%s|", node, comment)
			parser.popTailComment()
		}

		break
	}
}

func (parser *Parser) parseMethodDecl(contract *ast.Contract) bool {
	for parser.parseAnnotation() {

	}

	token := parser.peek()

	if token.Type != lexer.TokenType('}') {

		returnVal := parser.expectTypeDecl("expect method return type")

		tokenName := parser.expectf(lexer.TokenID, "expect method name")

		name := tokenName.Value.(string)

		parser.D("parse method %s", name)

		method, ok := contract.NewMethod(name)

		if !ok {
			parser.errorf(token.Start, "duplicate contract(%s) field(%s)", contract, name)
		}

		parser.attachAnnotation(method)

		method.Return = returnVal

		parser.parseParams(method)

		parser.parseExceptions(method)

		end := parser.expectf(lexer.TokenType(';'), "expect method name").End

		_setNodePos(method, token.Start, end)

		parser.attachComment(method)

		return true
	}

	return false
}

func (parser *Parser) parseExceptions(method *ast.Method) {

	token := parser.peek()

	if token.Type != lexer.KeyThrows {
		return
	}

	parser.next()

	parser.expectf(lexer.TokenType('('), "method exception table must start with (")

	for {

		typeDecl := parser.expectTypeDecl("expect exception type")

		exception := method.NewException(typeDecl)

		start, end := Pos(typeDecl)

		_setNodePos(exception, start, end)

		token := parser.peek()

		if token.Type != lexer.TokenType(',') {
			break
		}

		parser.next()
	}

	parser.expectf(lexer.TokenType(')'), "method exception table must end with )")
}

func (parser *Parser) parseParams(method *ast.Method) {

	parser.expectf(lexer.TokenType('('), "method param table must start with (")

	for {

		token := parser.peek()

		if token.Type == lexer.TokenType(')') {
			break
		}

		for parser.parseAnnotation() {

		}

		typeDecl := parser.expectTypeDecl("expect method param type declare")

		nameToken := parser.expectf(lexer.TokenID, "expect method param name")

		name := nameToken.Value.(string)

		param, ok := method.NewParam(name, typeDecl)

		parser.attachAnnotation(param)

		if !ok {
			parser.errorf(token.Start, "duplicate method(%s) param(%s)", method, name)
		}

		_setNodePos(param, token.Start, nameToken.End)
	}

	parser.expectf(lexer.TokenType(')'), "method param table must end with )")
}

func (parser *Parser) parseFieldDecl(table *ast.Table) bool {

	token := parser.peek()

	if token.Type != lexer.TokenType('}') {

		for parser.parseAnnotation() {

		}

		typeDecl := parser.expectTypeDecl("expect table(%s) field type declare", table)

		tokenName := parser.expectf(lexer.TokenID, "expect table(%s) field name", table)

		parser.expectf(lexer.TokenType(';'), "expect table(%s) field end tag ;", table)

		name := tokenName.Value.(string)

		field, ok := table.NewField(name, typeDecl)

		if !ok {
			parser.errorf(token.Start, "duplicate table(%s) field(%s)", table, name)
		}

		parser.attachAnnotation(field)

		_setNodePos(field, token.Start, tokenName.End)

		parser.attachComment(field)

		return true
	}

	return false
}

func (parser *Parser) expectTypeDecl(fmtstring string, args ...interface{}) (typeDecl ast.Type) {

	msg := fmt.Sprintf(fmtstring, args...)

	for {
		token := parser.peek()

		switch token.Type {
		case lexer.KeyByte, lexer.KeySByte, lexer.KeyInt16, lexer.KeyUInt16,
			lexer.KeyInt32, lexer.KeyUInt32, lexer.KeyInt64, lexer.KeyUInt64,
			lexer.KeyFloat32, lexer.KeyFloat64, lexer.KeyString, lexer.KeyBool, lexer.KeyVoid:

			typeDecl = ast.NewBuiltinType(token.Type)

			_setNodePos(typeDecl, token.Start, token.End)

			parser.next()

		case lexer.TokenID:
			name, star, end := parser.expectFullName("expect type declare")

			typeDecl = ast.NewTypeRef(name)

			_setNodePos(typeDecl, star, end)

		default:
			parser.errorf(token.Start, "%s\n\tunexpect token %s", msg, token)
			continue
		}

		for {

			if seqType, ok := parser.parseSeq(typeDecl); ok {
				typeDecl = seqType
				continue
			}

			break
		}

		return
	}

}

func (parser *Parser) parseSeq(component ast.Type) (typeDecl ast.Type, ok bool) {

	token := parser.peek()

	if token.Type != lexer.TokenType('[') {
		return nil, false
	}

	ok = true

	parser.next()

	token = parser.peek()

	if token.Type == lexer.TokenINT {
		parser.next()
		typeDecl = ast.NewSeq(component, int(token.Value.(int64)))
	} else {
		typeDecl = ast.NewSeq(component, -1)
	}

	end := parser.expectf(lexer.TokenType(']'), "seq type must end with ]").End

	start, _ := Pos(component)

	_setNodePos(typeDecl, start, end)

	return
}

func (parser *Parser) expectArgsTable(fmtstring string, args ...interface{}) (expr *ast.ArgsTable) {

	msg := fmt.Sprintf(fmtstring, args...)

	for {

		token := parser.peek()

		expr := parser.parseArgsTable()

		if expr != nil {
			return expr
		}

		parser.errorf(token.Start, msg)

		parser.next()
	}
}

func (parser *Parser) parseArgsTable() *ast.ArgsTable {

	token := parser.peek()

	if token.Type != lexer.TokenType('(') {

		return nil
	}

	parser.next()

	token = parser.peek()

	if token.Value == lexer.TokenType(')') {
		return ast.NewArgsTable(true)
	}

	token = parser.peek()

	start := token.Start

	end := token.End

	var args *ast.ArgsTable

	// this is named args table
	if token.Type == lexer.TokenLABEL {
		args = ast.NewArgsTable(true)

		for {

			token := parser.expectf(lexer.TokenLABEL, "expect arg label")

			label := token.Value.(string)

			arg := parser.expectArg("expect label(%s) value", label)

			parser.D("lable:%s", label)

			namedArg := ast.NewNamedArg(label, arg)

			_, end = Pos(arg)

			_setNodePos(namedArg, token.Start, end)

			args.Append(namedArg)

			if parser.peek().Type != lexer.TokenType(',') {
				break
			}
		}

	} else {

		args = ast.NewArgsTable(false)

		for {

			arg := parser.expectArg("expect arg")

			args.Append(arg)

			if parser.peek().Type != lexer.TokenType(',') {
				break
			}

			parser.next()
		}
	}

	_setNodePos(args, start, end)

	parser.expectf(lexer.TokenType(')'), "arg table must end with ')'")

	return args
}

func (parser *Parser) expectArg(fmtstring string, args ...interface{}) (expr ast.Expr) {
	msg := fmt.Sprintf(fmtstring, args...)

	for {

		token := parser.peek()

		expr = parser.parseArg()

		if expr != nil {
			return
		}

		parser.errorf(token.Start, msg)

		parser.next()
	}

}

func (parser *Parser) parseArg() ast.Expr {

	token := parser.peek()

	switch token.Type {
	case lexer.TokenINT, lexer.TokenFLOAT, lexer.TokenSTRING, lexer.TokenTrue, lexer.TokenFalse, lexer.TokenID:
		return parser.expectExpr("expect arg expr")
	case lexer.OpPlus:
		parser.next()
		numeric := parser.expectNumeric("unary op %s expect numeric object", lexer.OpPlus)

		_, end := Pos(numeric)

		_setNodePos(numeric, token.Start, end)

		return numeric
	case lexer.OpSub:
		parser.next()
		numeric := parser.expectNumeric("unary op %s expect numeric object", lexer.OpSub)

		numeric.Val = -numeric.Val

		_, end := Pos(numeric)

		_setNodePos(numeric, token.Start, end)

		return numeric
	default:

		return nil
	}
}

func (parser *Parser) expectExpr(fmtStr string, args ...interface{}) ast.Expr {

	msg := fmt.Sprintf(fmtStr, args...)

	for {
		token := parser.peek()

		var expr ast.Expr

		switch token.Type {

		case lexer.TokenINT:
			parser.next()

			expr = ast.NewNumeric(float64(token.Value.(int64)))

			_setNodePos(expr, token.Start, token.End)

		case lexer.TokenFLOAT:

			parser.next()

			expr = ast.NewNumeric(token.Value.(float64))

			_setNodePos(expr, token.Start, token.End)

		case lexer.TokenSTRING:
			parser.next()

			expr = ast.NewString(token.Value.(string))

			_setNodePos(expr, token.Start, token.End)

		case lexer.TokenTrue:
			parser.next()

			expr = ast.NewBoolean(true)

			_setNodePos(expr, token.Start, token.End)

		case lexer.TokenFalse:
			parser.next()

			expr = ast.NewBoolean(false)

			_setNodePos(expr, token.Start, token.End)

		case lexer.TokenID:
			name, start, end := parser.expectFullName("expect constant reference or table instance")

			token = parser.peek()

			if token.Type == lexer.TokenType('(') {
				initargs := parser.expectArgsTable("expect table instance init args table")
				newObj := ast.NewNewObj(name, initargs)

				_setNodePos(newObj, start, end)

				return newObj
			}

			expr = ast.NewConstantRef(name)

			_setNodePos(expr, start, end)
		}

		if expr != nil {

			token := parser.peek()

			switch token.Type {
			case lexer.OpBitOr, lexer.OpBitAnd:

				parser.next()

				rhs := parser.expectExpr("expect binary op(%s) rhs", token.Type)

				binaryOp := ast.NewBinaryOp(token.Type, expr, rhs)

				start, _ := Pos(expr)

				_, end := Pos(binaryOp)

				_setNodePos(binaryOp, start, end)

				expr = binaryOp
			}

			return expr
		}

		parser.errorf(token.Start, msg)

		parser.next()
	}
}

func (parser *Parser) expectNumeric(fmtStr string, args ...interface{}) (number *ast.Numeric) {
	msg := fmt.Sprintf(fmtStr, args...)

	for {
		token := parser.next()

		parser.D("expect numeric :%s", token)

		if token.Type == lexer.TokenINT {
			number = ast.NewNumeric(float64(token.Value.(int64)))
		} else if token.Type == lexer.TokenFLOAT {
			number = ast.NewNumeric(token.Value.(float64))
		}

		if number != nil {
			return number
		}

		parser.errorf(token.Start, msg)
	}

}

func (parser *Parser) expectFullName(fmtstring string, args ...interface{}) (string, lexer.Position, lexer.Position) {
	msg := fmt.Sprintf(fmtstring, args...)

	var buff bytes.Buffer

	token := parser.expectf(lexer.TokenID, msg)

	buff.WriteString(token.Value.(string))

	start := token.Start

	end := token.End

	for {
		token = parser.peek()

		if token.Type != lexer.TokenType('.') {
			break
		}

		buff.WriteRune('.')

		parser.next()

		token = parser.expectf(lexer.TokenID, msg)

		buff.WriteString(token.Value.(string))

		end = token.End
	}

	return buff.String(), start, end
}

func (parser *Parser) parseAnnotation() bool {

	for parser.parseComment() {
	}

	token := parser.peek()

	if token.Type != lexer.TokenType('@') {
		return false
	}

	start := token.Start

	parser.next()

	name, start, end := parser.expectFullName("expect annotation name")

	annotation := ast.NewAnnotation(name)

	parser.D("annotation [%s]", name)

	token = parser.peek()

	if token.Type == lexer.TokenType('(') {

		args := parser.expectArgsTable("expect annotation arg table")

		_, end = Pos(args)

		annotation.Args = args
	}

	_setNodePos(annotation, start, end)

	parser.annotationStack = append(parser.annotationStack, annotation)

	return true
}

func (parser *Parser) parseComment() bool {

	token := parser.peek()
	if token.Type != lexer.TokenCOMMENT {
		return false
	}

	//move to next token
	parser.next()

	if len(parser.commentStack) > 0 {

		comment := parser.commentStack[len(parser.commentStack)-1]

		var pos lexer.Position

		val, ok := comment.GetExtra("end")

		if ok {
			pos = val.(lexer.Position)
		}

		if pos.Lines+1 == token.Start.Lines {
			comment.SetExtra("end", token.End)

			comment.Append(token.Value.(string))

			return true
		}
	}

	comment := ast.NewComment()

	comment.SetExtra("start", token.Start)
	comment.SetExtra("end", token.End)

	comment.Append(token.Value.(string))

	parser.commentStack = append(parser.commentStack, comment)

	return true
}

func (parser *Parser) tailComment() (*ast.Comment, bool) {
	if len(parser.commentStack) == 0 {
		return nil, false
	}

	return parser.commentStack[len(parser.commentStack)-1], true
}

func (parser *Parser) popTailComment() (*ast.Comment, bool) {
	if len(parser.commentStack) == 0 {
		return nil, false
	}

	comment := parser.commentStack[len(parser.commentStack)-1]

	parser.commentStack = parser.commentStack[:len(parser.commentStack)-1]

	return comment, true
}

func (parser *Parser) parseImport() bool {

	for parser.parseComment() {
	}

	token := parser.peek()

	//the import instructions must be typed at the beginning of script
	if token.Type != lexer.KeyImport {
		return false
	}

	// move next token
	parser.next()

	usingNamePath, start, end := parser.expectFullName("expect using name path")

	using := parser.script.Using(usingNamePath)

	parser.expectf(lexer.TokenType(';'), "import name path must end with ';'")

	parser.D("parse using :%s", using)

	_setNodePos(using, start, end)

	parser.attachComment(using)

	return true
}
