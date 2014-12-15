package gslang

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"strings"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslogger"
)

//Parser's public error codes
var (
	ErrParse = errors.New("gslang parser error")
)

const (
	posExtra     = "gslang_parser_pos"
	commentExtra = "gslang_parser_comment"
)

func attachPos(node ast.Node, pos Position) {
	node.NewExtra(posExtra, pos)
}

//Pos get the AST node's pos extra data
func Pos(node ast.Node) Position {
	if val, ok := node.Extra(posExtra); ok {
		return val.(Position)
	}

	return Position{
		FileName: "<unknown>",
		Lines:    0,
		Column:   0,
	}
}

func attachComments(node ast.Node, comments []*Token) {
	node.NewExtra(commentExtra, comments)
}

//Comments get the AST node's comments extra data
func Comments(node ast.Node) []*Token {
	if val, ok := node.Extra(commentExtra); ok {
		return val.([]*Token)
	}

	return nil
}

type commentStack []*Token

func (stack *commentStack) push(token *Token) {
	gserrors.Require(token.Type == TokenCOMMENT, "require push comment token")
	*stack = append(*stack, token)
}

func (stack *commentStack) pop() (token *Token, ok bool) {
	if len(*stack) == 0 {
		return nil, false
	}

	token = (*stack)[len(*stack)-1]

	*stack = (*stack)[:len(*stack)-1]

	ok = true

	return
}

//Parser gslang script parser
type Parser struct {
	gslogger.Log             //Mixin log implement
	*Lexer                   //Mixin lexer implement
	cs           *CompileS   //parser belongs compile service object
	script       *ast.Script //AST script node object which the parser try to build
	comments     []*Token    //cached Q for comments
	attrs        []*ast.Attr //cached Q for attrs
}

// Peek override lexer's Peek function :
// parser's Peek function will panic when error occur
func (parser *Parser) Peek() *Token {
	token, err := parser.Lexer.Peek()
	if err != nil {
		panic(err)
	}

	return token
}

// Next override lexer's Next function:
// parser's Next function will panic when error occur
func (parser *Parser) Next() (token *Token) {
	token, err := parser.Lexer.Next()
	if err != nil {
		panic(err)
	}

	return token
}

func (parser *Parser) errorf(position Position, fmtstring string, args ...interface{}) {
	gserrors.Panicf(
		ErrParse,
		fmt.Sprintf(
			"parse %s error : %s",
			position,
			fmt.Sprintf(fmtstring, args...),
		),
	)
}

func (parser *Parser) errorf2(err error, position Position, fmtstring string, args ...interface{}) {
	gserrors.Panicf(
		err,
		fmt.Sprintf(
			"parse %s error : %s",
			position,
			fmt.Sprintf(fmtstring, args...),
		),
	)
}

func (parser *Parser) expect(expect rune) *Token {
	token := parser.Next()
	if token.Type != expect {
		parser.errorf(token.Pos, "expect '%s',but got '%s'", TokenName(expect), TokenName(token.Type))
	}

	return token
}

func (parser *Parser) expectf(expect rune, fmtstring string, args ...interface{}) *Token {
	token := parser.Next()
	if token.Type != expect {
		parser.errorf(token.Pos, fmt.Sprintf(fmtstring, args...))
	}

	return token
}

func (parser *Parser) parseTypeRef() *ast.TypeRef {
	start := parser.expect(TokenID)

	nodes := []string{start.Value.(string)}

	for {
		token := parser.Peek()

		if token.Type != '.' {
			break
		}
		parser.Next()
		token = parser.expect(TokenID)

		nodes = append(nodes, token.Value.(string))
	}

	ref := parser.script.NewTypeRef(nodes)

	attachPos(ref, start.Pos)

	return ref
}

func (cs *CompileS) parse(pkg *ast.Package, path string) (*ast.Script, error) {

	script, err := pkg.NewScript(filepath.Base(path))

	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	parser := &Parser{
		Log:    gslogger.Get("gslang[parser]"),
		Lexer:  NewLexer(script.Name(), bytes.NewBuffer(content)),
		cs:     cs,
		script: script,
	}

	err = parser.parse()

	return script, err
}

func (parser *Parser) parse() (err error) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(gserrors.GSError); ok {
				err = e.(error)
			} else {
				err = gserrors.New(e.(error))
			}

		}
	}()
	parser.parseImports()

	for {

		parser.parseAttrs()
		token := parser.Next() //lookup next token
		switch token.Type {
		case TokenEOF:
			goto FINISH
		case KeyEnum:
			parser.parseEnum()
		case KeyTable:
			parser.parseTable(false)
		case KeyStruct:
			parser.parseTable(true)
		case KeyContract:
			parser.parseContract()
		default:
			parser.errorf(token.Pos, "expect EOF")
		}
	}
FINISH:
	//attach rest of comments to script node
	attachComments(parser.script, parser.comments)
	parser.script.AddAttrs(parser.attrs)
	return
}

func (parser *Parser) parseContract() {
	name := parser.expect(TokenID)

	contract := parser.script.NewContract(name.Value.(string))

	if old, ok := parser.script.NewType(contract); !ok {
		parser.errorf(name.Pos, "duplicate type name :\n\tsee:%s", Pos(old))
	}

	attachPos(contract, name.Pos)
	parser.attachComments(contract)
	parser.attachAttrs(contract)

	token := parser.Peek()

	if token.Type == '(' { //parse inher table
		parser.Next()

		for {
			parser.parseComments()
			base := parser.parseTypeRef()

			if old, ok := contract.NewBase(base); ok {
				parser.parseComments()
				parser.attachComments(base)
			} else {
				parser.errorf(
					Pos(base),
					"duplicate inher from same contract :\n\tsee:%s",
					Pos(old),
				)
			}

			next := parser.Peek()

			if next.Type == ',' {
				parser.Next()
				continue
			}

			break //break parse inher table
		}

		parser.expect(')')
	}

	parser.expect('{')

	for {

		parser.parseAttrs()

		token := parser.Peek()

		if token.Type != TokenID {
			break
		}

		methodName := parser.Next()

		method, ok := contract.NewMethod(methodName.Value.(string))

		if !ok {
			parser.errorf(methodName.Pos, "duplicate method name :\n\tsee:%s", Pos(method))
		}

		attachPos(method, methodName.Pos)
		parser.attachComments(method)
		parser.attachAttrs(method)

		parser.expect('(')

		next := parser.Peek()

		if next.Type != ')' { //parse params list

			for {

				parser.parseAttrs()

				parmaType := parser.parseType()

				next := parser.Peek()

				if next.Type != ',' &&
					next.Type != ')' &&
					next.Type != TokenCOMMENT {
					parmaType = parser.parseType()
				}

				param := method.NewParam(parmaType)
				attachPos(param, Pos(param.Type))

				//attach comments and attrs
				parser.parseComments()
				parser.attachComments(param)
				parser.attachAttrs(param)
				next = parser.Peek()

				if next.Type == ',' {
					parser.Next()
					continue
				}

				break
			}

		}

		parser.expect(')')

		next = parser.Peek()

		if next.Type == TokenArrowRight { //parse return params list
			parser.Next()
			parser.expect('(')
			for {

				parser.parseAttrs()
				parmaType := parser.parseType()

				next := parser.Peek()

				if next.Type != ',' &&
					next.Type != ')' &&
					next.Type != TokenCOMMENT {
					parmaType = parser.parseType()
				}

				param := method.NewReturn(parmaType)
				attachPos(param, Pos(param.Type))

				//attach comments and attrs
				parser.parseComments()
				parser.attachComments(param)
				parser.attachAttrs(param)

				next = parser.Peek()

				if next.Type == ',' {
					parser.Next()
					continue
				}

				break
			}
			parser.expect(')')
		}

		parser.expect(';')
	}

	parser.expect('}')
}

func (parser *Parser) newGsLangAttr(name string) *ast.Attr {

	if parser.script.Package().Name() != GSLangPackage {

		return parser.script.NewAttr(
			parser.script.NewTypeRef([]string{
				"gslang", name,
			},
			),
		)
	}

	return parser.script.NewAttr(parser.script.NewTypeRef([]string{name}))
}

func (parser *Parser) newGsLangTypeRef(name string) *ast.TypeRef {
	if parser.script.Package().Name() != GSLangPackage {
		return parser.script.NewTypeRef([]string{
			"gslang", name,
		})
	}

	return parser.script.NewTypeRef([]string{name})
}

// parseType parse type expression
func (parser *Parser) parseType() ast.Expr {
	token := parser.Peek()

	switch token.Type {
	case '[':
		parser.Next()
		next := parser.Peek()
		length := uint16(0)
		if next.Type == TokenINT {
			parser.Next()
			//check array length range
			val := next.Value.(int64)
			if val < 1 || val > math.MaxUint16 {
				parser.errorf(next.Pos, "array lenght out of range :%d", val)
			}

			length = uint16(val)
		}

		parser.expect(']')

		component := parser.parseType()

		//check Recursively array/list define
		switch component.(type) {
		case *ast.List, *ast.Array:
			parser.errorf(token.Pos, "gslang didn't support Recursively define array or list")
		}

		var expr ast.Expr
		if length > 0 { //this is array
			expr = parser.script.NewArray(length, component)
		} else {
			expr = parser.script.NewList(component)
		}

		attachPos(expr, token.Pos)

		return expr

	case
		KeyByte, KeySByte, KeyInt16, KeyUInt16, KeyInt32, KeyUInt32,
		KeyInt64, KeyUInt64, KeyBool, KeyFloat32, KeyFloat64, KeyString:
		parser.Next()
		expr := parser.newGsLangTypeRef(strings.Title(TokenName(token.Type)))
		attachPos(expr, token.Pos)
		return expr
	case TokenID:
		expr := parser.parseTypeRef()
		attachPos(expr, token.Pos)
		return expr
	default:
		parser.errorf(token.Pos, "expect type declare")
	}

	return nil
}

// parseTable parse table type
func (parser *Parser) parseTable(isStruct bool) {

	name := parser.expect(TokenID)
	table := parser.script.NewTable(name.Value.(string))

	if old, ok := parser.script.NewType(table); !ok {
		parser.errorf(name.Pos, "duplicate type name :\n\tsee:%s", Pos(old))
	}

	attachPos(table, name.Pos)
	parser.attachComments(table)

	parser.attachAttrs(table)

	if isStruct {
		attr := parser.newGsLangAttr("Struct")
		attachPos(attr, name.Pos)
		table.AddAttr(attr)
	}

	parser.expect('{')

	for {
		parser.parseAttrs()

		token := parser.Peek()

		if token.Type != TokenID {
			break
		}

		fieldName := parser.expect(TokenID)
		field, ok := table.NewField(fieldName.Value.(string))

		if !ok {
			parser.errorf(fieldName.Pos, "duplicate field name:\n\t see:%s", Pos(field))
		}

		attachPos(field, fieldName.Pos)
		field.Type = parser.parseType()
		parser.expect(';')
		parser.parseComments()
		parser.attachComments(field)
		parser.attachAttrs(field)
	}

	parser.expect('}')
}

// parseEnumBase parse enum inherit type to calc enum type len
func (parser *Parser) parseEnumBase() (length uint, signed bool) {
	parser.expect('(')
	token := parser.Next()

	switch token.Type {
	case KeyByte:
		length = 1
	case KeySByte:
		length = 1
		signed = true
	case KeyInt16:
		length = 2
		signed = true
	case KeyUInt16:
		length = 2
		signed = false
	case KeyInt32:
		length = 4
		signed = true
	case KeyUInt32:
		length = 4
		signed = false
	default:
		parser.errorf(
			token.Pos,
			"enum must inherit from integer types , got :%s",
			TokenName(token.Type),
		)
	}

	parser.expect(')')

	return
}

// parseEnum parse the enum type
func (parser *Parser) parseEnum() {

	name := parser.expect(TokenID)

	token := parser.Peek()

	length := uint(1)

	signed := false

	if token.Type == '(' {
		length, signed = parser.parseEnumBase()
	}

	enum := parser.script.NewEnum(name.Value.(string), length, signed)

	if old, ok := parser.script.NewType(enum); !ok {
		parser.errorf(name.Pos, "duplicate type name :\n\tsee:%s", Pos(old))
	}

	attachPos(enum, name.Pos)
	parser.attachComments(enum)
	parser.attachAttrs(enum)

	parser.expect('{')

	for {
		parser.parseAttrs()

		token := parser.expectf(TokenID, "expect enum value field")
		//parse enum val field
		parser.expect('(')

		next := parser.Peek()

		negative := false

		if next.Type == '-' {
			parser.Next()
			negative = true
		}

		valToken := parser.expect(TokenINT)

		val := valToken.Value.(int64)

		if negative {
			val = -val
		}

		//check the enum val range
		switch {
		case enum.Length == 1 && enum.Signed == true:
			if val > math.MaxInt8 || val < math.MinInt8 {
				parser.errorf(valToken.Pos, "out of enum[%s] type's range", enum)
			}
		case enum.Length == 1 && enum.Signed == false:
			if val > math.MaxUint8 || val < 0 {
				parser.errorf(valToken.Pos, "out of enum[%s] type's range", enum)
			}
		case enum.Length == 2 && enum.Signed == true:
			if val > math.MaxInt16 || val < math.MinInt16 {
				parser.errorf(valToken.Pos, "out of enum[%s] type's range", enum)
			}
		case enum.Length == 2 && enum.Signed == false:
			if val > math.MaxUint16 || val < 0 {
				parser.errorf(valToken.Pos, "out of enum[%s] type's range", enum)
			}
		case enum.Length == 4 && enum.Signed == true:
			if val > math.MaxInt32 || val < math.MinInt32 {
				parser.errorf(valToken.Pos, "out of enum[%s] type's range", enum)
			}
		case enum.Length == 4 && enum.Signed == false:
			if val > math.MaxUint32 || val < 0 {
				parser.errorf(valToken.Pos, "out of enum[%s] type's range", enum)
			}
		}

		parser.expect(')')

		enumVal, ok := enum.NewVal(token.Value.(string), val)

		if !ok {
			parser.errorf(
				token.Pos,
				"duplicate enum val name(%s):\n\tsee:%s",
				enumVal.Name(),
				Pos(enumVal),
			)
		}

		attachPos(enumVal, token.Pos) //bind position extra data
		parser.attachAttrs(enumVal)   //bind attrs

		// if not found ','  break field parse loop
		next = parser.Peek()
		if next.Type != ',' {
			parser.parseComments()
			parser.attachComments(enumVal)
			break
		}
		parser.Next()
		parser.parseComments()
		parser.attachComments(enumVal)
	}

	parser.expect('}')
}

func (parser *Parser) parseAttrs() {
	parser.parseComments()

	for {

		token := parser.Peek()
		if token.Type != '@' {
			return
		}
		//move cousor to next token
		parser.Next()

		//create new attr node obj
		attr := parser.script.NewAttr(
			parser.parseTypeRef(),
		)
		//bind pos to attr node
		attachPos(attr, token.Pos)
		//check if this attr has argument list
		token = parser.Peek()
		//parse attribute arguments
		if token.Type == '(' {
			parser.Next()
			attr.Args = parser.parseArgs()
			parser.expect(')') //expect argument list end token
		}

		//parser.expect(']') //parse attribute end token

		parser.attrs = append(parser.attrs, attr) //add attr to canched Q

		parser.parseComments() //parse comment same line with attr

		parser.attachComments(attr) //attach comments
	}

}

// parseArgs parse the arg list include named argument list
func (parser *Parser) parseArgs() ast.Expr {
	//first check if the argument list is named argument list
	token := parser.Peek()

	if token.Type == ')' {
		return nil //empty argument list
	}

	if token.Type == TokenLABEL {
		//parse named argument list

		args := parser.script.NewNamedArgs()

		parser.Next()

		name := token

		for {
			if arg, ok := args.NewArg(token.Value.(string), parser.parseArg()); !ok {
				parser.errorf(
					name.Pos,
					"duplicate param assign :\n\t see :%s",
					Pos(arg),
				)
			} else {
				parser.parseComments()
				parser.attachComments(arg)
			}

			token = parser.Peek()

			if token.Type != ',' {
				break //no more argments
			}

			parser.Next()

			name = parser.expect(TokenLABEL)
		}

		return args
	}

	args := parser.script.NewArgs()

	for {
		arg := args.NewArg(parser.parseArg())
		parser.parseComments()
		parser.attachComments(arg)
		token = parser.Peek()

		if token.Type != ',' {
			break
		}
	}

	return args
}

// parseArg parse one argument stmt
func (parser *Parser) parseArg() ast.Expr {
	var lhs *ast.BinaryOp

	for {
		token := parser.Peek()
		var rhs ast.Expr
		switch token.Type {
		case TokenINT:
			parser.Next()
			rhs = parser.script.NewInt(token.Value.(int64))
		case TokenFLOAT:
			parser.Next()
			rhs = parser.script.NewFloat(token.Value.(float64))
		case TokenSTRING:
			parser.Next()
			rhs = parser.script.NewString(token.Value.(string))
		case TokenTrue:
			parser.Next()
			rhs = parser.script.NewBool(true)
		case TokenFalse:
			parser.Next()
			rhs = parser.script.NewBool(false)
		case '-':
			parser.Next()
			next := parser.Next()

			if next.Type == TokenINT {
				rhs = parser.script.NewInt(-next.Value.(int64))
			} else if next.Type == TokenFLOAT {
				rhs = parser.script.NewFloat(-next.Value.(float64))
			}

			parser.errorf(token.Pos, "unexpect token '-'")

		case '+':
			parser.Next()
			next := parser.Next()

			if next.Type == TokenINT {
				rhs = parser.script.NewInt(next.Value.(int64))
			} else if next.Type == TokenFLOAT {
				rhs = parser.script.NewFloat(next.Value.(float64))
			}
			parser.errorf(token.Pos, "unexpect token '+'")
		case TokenID:
			rhs = parser.parseTypeRef()
		default:
			parser.errorf(token.Pos, "unexpect token '%s', expect argument stmt", TokenName(token.Type))
		}

		attachPos(rhs, token.Pos)

		if lhs != nil {
			lhs.Right = rhs
		}

		token = parser.Peek()

		if token.Type == '|' {
			parser.Next()

			if lhs != nil {
				lhs = parser.script.NewBinaryOp("|", lhs, nil)
			} else {
				lhs = parser.script.NewBinaryOp("|", rhs, nil)
			}

			attachPos(lhs, token.Pos)

			continue
		}

		if lhs != nil {
			return lhs
		}

		return rhs
	}
}

// parseComment parse comment tokens may occur, and cached
func (parser *Parser) parseComments() {

	for {
		token := parser.Peek()
		if token.Type != TokenCOMMENT {
			return
		}

		//move to next token
		parser.Next()
		//cache the comment token
		parser.comments = append(parser.comments, token)
	}
}

func (parser *Parser) attachComments(node ast.Node) {
	pos := Pos(node)
	gserrors.Assert(
		pos.Valid(),
		"all node have to bind with pos object by calling attachPos :%s",
		node,
	)

	var selected []*Token

	var rest []*Token

	for i := len(parser.comments) - 1; i >= 0; i-- {
		comment := parser.comments[i]
		gserrors.Assert(
			pos.FileName == comment.Pos.FileName,
			"comment's filename must equal with node's one",
		)

		if comment.Pos.Lines == pos.Lines ||
			(comment.Pos.Lines+1) == pos.Lines {

			selected = append(selected, comment)
			pos = comment.Pos //set the new pos
		} else {
			rest = append(rest, comment)
		}
	}

	parser.comments = rest

	var revert []*Token

	for i := len(selected) - 1; i >= 0; i-- {
		revert = append(revert, selected[i])
	}

	attachComments(node, revert)
}

func (parser *Parser) attachAttrs(node ast.Node) {
	node.AddAttrs(parser.attrs)
	parser.attrs = nil
}

// parseImports parse the import instruction at the beginning of this script
func (parser *Parser) parseImports() {
	//cache possible comment token
	parser.parseComments()
	for {
		token := parser.Peek()
		//the import instructions must be typed at the beginning of script
		if token.Type != KeyImport {
			break
		}
		parser.Next()
		//parse the import body :maybe TokenSTRING or starting with '('
		token = parser.Peek()

		if token.Type == TokenSTRING || token.Type == TokenID {
			ref := parser.parseImport()
			gserrors.Assert(ref != nil, "check parser.parseImport implement")
			parser.parseComments()
			parser.attachComments(ref)
		} else if token.Type == '(' {
			parser.Next()
			//parse import instructions between '()'
			for {
				parser.parseComments()
				if ref := parser.parseImport(); ref != nil {
					parser.parseComments()
					parser.attachComments(ref)
					continue
				}

				break
			}
			parser.expect(')')
		} else {
			parser.errorf(token.Pos, "expect import body: TokenSTRING or '('")
		}
	}

	// all script have to import package "github.com/gsdocker/gslang"
	// the parser auto add it into every script's import table except
	// scripts in package "github.com/gsdocker/gslang"
	if parser.script.Package().Name() != GSLangPackage &&
		parser.script.Imports["gslang"] == nil {

		pkg, err := parser.cs.Compile(GSLangPackage)

		if err != nil {
			panic(err)
		}

		pos := Position{
			FileName: parser.script.Name(),
			Lines:    1,
			Column:   1,
		}

		ref, ok := parser.script.NewPackageRef("gslang", pkg)

		gserrors.Assert(pkg != nil, "check CompileS#Compile implement")

		gserrors.Assert(ok, "must check : if the script manual import gslang package")

		attachPos(ref, pos)
	}

}

func (parser *Parser) parseImport() *ast.PackageRef {
	token := parser.Peek()
	var (
		path string
		key  string
	)
	if token.Type == TokenSTRING {
		path = token.Value.(string)
		key = filepath.Base(path)
		parser.Next()
	} else if token.Type == TokenID {
		parser.Next()
		key = token.Value.(string) //get key
		token := parser.expect(TokenSTRING)
		path = token.Value.(string) //get path
	} else {
		return nil
	}

	//load import package
	pkg, err := parser.cs.Compile(path)

	if err != nil {
		panic(err)
	}

	ref, ok := parser.script.NewPackageRef(key, pkg)

	if !ok {
		parser.errorf(token.Pos, "import same package(%s) twice : \n\tsee :%s", key, Pos(ref))
	}

	attachPos(ref, token.Pos)

	return ref
}
