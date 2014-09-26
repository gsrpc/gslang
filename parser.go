package gslang

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

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
	gslogger.Log              //Mixin log implement
	*Lexer                    //Mixin lexer implement
	cs           *CompileS    //parser belongs compile service object
	script       *ast.Script  //AST script node object which the parser try to build
	comments     commentStack //canced comments
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
	gserrors.Newf(
		ErrParse,
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
			err = e.(error)
		}
	}()
	parser.parseImports()

	for {
		token := parser.Peek() //lookup next token

		switch token.Type {
		case TokenEOF:
			goto FINISH
		case KeyEnum:

		default:
			parser.errorf(token.Pos, "expect EOF")
		}
	}
FINISH:
	//attach rest of comments to script node
	attachComments(parser.script, parser.comments)

	return
}

func (parser *Parser) parseAttr() {
	parser.parseComments()

	for {
		token := parser.Peek()
		if token.Type != '[' {
			return
		}
		//move cousor to next token
		parser.Next()
		//create new attr node obj
		attr := parser.script.NewAttr(
			parser.expect(TokenID).Value.(string),
		)
		//bind pos to attr node
		attachPos(attr, token.Pos)
		//check if this attr has argument list
		token = parser.Peek()
		//parse attribute arguments
		if token.Type == '(' {
			//TODO: add args parse codes
			parser.expect(')') //expect argument list end token
		}

		parser.expect(']') //parse attribute end token
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
		parser.comments.push(token)
	}
}

func (parser *Parser) attachComments(node ast.Node) {
	pos := Pos(node)
	gserrors.Assert(
		pos.Valid(),
		"all node have to bind with pos object by calling attachPos",
	)

	var selected []*Token

	for {
		if comment, ok := parser.comments.pop(); ok {
			gserrors.Assert(
				pos.FileName == comment.Pos.FileName,
				"comment's filename must equal with node's one",
			)

			if comment.Pos.Lines == pos.Lines ||
				(comment.Pos.Lines+1) == pos.Lines {

				selected = append(selected, comment)
				pos = comment.Pos //set the new pos
				continue
			} else {
				parser.comments.push(comment)
			}
		}

		break
	}

	if selected != nil {
		parser.D("attach comments to node :%s", node)
		for _, token := range selected {
			parser.D("\t%s", token.Value.(string))
		}
	}

	attachComments(node, selected)
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
		parser.script.Imports[GSLangPackage] == nil {
		pkg, err := parser.cs.Compile(GSLangPackage)
		pos := Position{
			FileName: parser.script.Name(),
			Lines:    1,
			Column:   1,
		}
		if err != nil {
			parser.errorf(pos, "can't load import package\n\t:package:%s", GSLangPackage)
		}

		ref, ok := parser.script.NewPackageRef("gslang", pkg)

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
		parser.errorf(token.Pos, "can't load import package\n\t:package:%s", path)
	}

	ref, ok := parser.script.NewPackageRef(key, pkg)

	if !ok {
		parser.errorf(token.Pos, "import same package(%s) twice : \n\tsee :%s", key, Pos(ref))
	}

	attachPos(ref, token.Pos)

	return ref
}
