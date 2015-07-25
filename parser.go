package gslang

import (
	"bytes"
	"fmt"

	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslogger"
)

// Parser gslang parser
type Parser struct {
	gslogger.Log                      // mixin logger
	lexer           *Lexer            // token lexer
	script          *ast.Script       // current script
	commentStack    []*ast.Comment    // comment stack
	annotationStack []*ast.Annotation // comment stack
	errorHandler    ErrorHandler      // error Handler
}

func (compiler *Compiler) parse(lexer *Lexer, errorHandler ErrorHandler) *ast.Script {

	return (&Parser{
		Log:          gslogger.Get("parser"),
		lexer:        lexer,
		script:       ast.NewScript(lexer.String()),
		errorHandler: errorHandler,
	}).parse()
}

func (parser *Parser) peek() *Token {
	token, err := parser.lexer.Peek()
	if err != nil {
		panic(err)
	}

	return token
}

func (parser *Parser) next() (token *Token) {
	token, err := parser.lexer.Next()
	if err != nil {
		panic(err)
	}

	return token
}

func (parser *Parser) errorf(position Position, fmtstring string, args ...interface{}) {
	parser.errorHandler.HandleParseError(
		ErrParser,
		position,
		fmt.Sprintf(fmtstring, args...),
	)
}

func (parser *Parser) errorf2(err error, position Position, fmtstring string, args ...interface{}) {
	parser.errorHandler.HandleParseError(
		err,
		position,
		fmt.Sprintf(fmtstring, args...),
	)
}

func (parser *Parser) expectf(expect TokenType, fmtstring string, args ...interface{}) *Token {

	for {
		token := parser.next()

		if token.Type != expect {
			parser.errorf(token.Start, fmt.Sprintf(fmtstring, args...))
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

	for parser.parseAnnotation() {

	}

	return parser.script
}

func (parser *Parser) debug(fmtString string, args ...interface{}) {
	parser.Log.D("%s\n\tfile :%s", fmt.Sprintf(fmtString, args...), parser.script)
}

func (parser *Parser) parsePackage() {

	parser.debug("[] parse script's package line")

	for parser.parseComment() {
	}

	parser.expectf(KeyPackage, "script must start with package keyword")

	parser.script.Package, _, _ = parser.expectFullName("expect script's package name")

	parser.expectf(TokenType(';'), "package name must end with ';'")

	parser.debug("package [%s]", parser.script.Package)

	parser.debug("parse script's package line -- success")
}

func (parser *Parser) expectArgsTable(fmtstring string, args ...interface{}) (*ast.Expr, Position, Position) {
	return nil, Position{}, Position{}
}

func (parser *Parser) expectFullName(fmtstring string, args ...interface{}) (string, Position, Position) {
	msg := fmt.Sprintf(fmtstring, args...)

	var (
		buff  bytes.Buffer
		start Position
		end   Position
	)

	token := parser.expectf(TokenID, msg)

	buff.WriteString(token.Value.(string))

	start = token.Start

	for {
		token = parser.peek()

		if token.Type != TokenType('.') {
			break
		}

		buff.WriteRune('.')

		parser.next()

		token = parser.expectf(TokenID, msg)

		buff.WriteString(token.Value.(string))

		end = token.End
	}

	return buff.String(), start, end
}

func (parser *Parser) parseAnnotation() bool {

	for parser.parseComment() {
	}

	token := parser.peek()

	if token.Type != TokenType('@') {
		return false
	}

	parser.next()

	name, start, end := parser.expectFullName("expect annotation name")

	annotation := ast.NewAnnotation(name)

	parser.debug("annotation [%s]", name)

	token = parser.peek()

	if token.Type == TokenType('(') {
		parser.expectArgsTable("expect annotation arg table")
	}

	_setNodePos(annotation, start, end)

	parser.annotationStack = append(parser.annotationStack, annotation)

	return true
}

func (parser *Parser) parseComment() bool {

	token := parser.peek()
	if token.Type != TokenCOMMENT {
		return false
	}

	//move to next token
	parser.next()

	if len(parser.commentStack) > 0 {

		comment := parser.commentStack[len(parser.commentStack)-1]

		var pos Position

		comment.GetExtra("end", &pos)

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
	if token.Type != KeyImport {
		return false
	}

	// move next token
	parser.next()

	usingNamePath, start, end := parser.expectFullName("expect using name path")

	using := parser.script.Using(usingNamePath)

	_setNodePos(using, start, end)

	if comment, ok := parser.tailComment(); ok {
		if _AttachComment(using, comment) {
			parser.debug("attach using name path comment :\n|%s|", comment)
			parser.popTailComment()
		}
	}

	parser.debug("using [%s]", using)

	parser.expectf(TokenType(';'), "import name path must end with ';'")

	for parser.parseComment() {
		comment, _ := parser.tailComment()
		if _AttachComment(using, comment) {
			parser.debug("attach using name path comment :\n|%s|", comment)
			parser.popTailComment()
		}

		break
	}

	return true
}
