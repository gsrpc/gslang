package lexer

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslogger"
)

// Err
var ErrLexer = errors.New("gslang lexer error")

// TokenType type
type TokenType rune

//Token types
const (
	TokenEOF TokenType = -(iota + 1)
	TokenID
	TokenINT
	TokenFLOAT
	TokenTrue
	TokenFalse
	TokenSTRING
	TokenCOMMENT
	TokenLABEL
	TokenArrowRight
	KeyByte
	KeySByte
	KeyInt16
	KeyUInt16
	KeyInt32
	KeyUInt32
	KeyInt64
	KeyUInt64
	KeyFloat32
	KeyFloat64
	KeyString
	KeyBool
	KeyEnum
	KeyStruct
	KeyTable
	KeyContract
	KeyImport
	KeyPackage
	KeyVoid
	KeyThrows
	KeyType
	KeyMap
	OpBitOr
	OpBitAnd
	OpPlus
	OpSub
)

var tokenName = map[TokenType]string{
	TokenEOF:        "EOF",
	TokenID:         "ID",
	TokenINT:        "INT",
	TokenFLOAT:      "FLOAT",
	TokenSTRING:     "STRING",
	TokenCOMMENT:    "COMMENT",
	TokenLABEL:      "LABLE",
	TokenArrowRight: "->",
	KeyByte:         "byte",
	KeySByte:        "sbyte",
	KeyInt16:        "int16",
	KeyUInt16:       "uint16",
	KeyInt32:        "int32",
	KeyUInt32:       "uint32",
	KeyInt64:        "int64",
	KeyUInt64:       "uint64",
	KeyFloat32:      "float32",
	KeyFloat64:      "float64",
	KeyString:       "string",
	KeyBool:         "bool",
	KeyEnum:         "enum",
	KeyStruct:       "struct",
	KeyTable:        "table",
	KeyContract:     "contract",
	KeyImport:       "using",
	KeyPackage:      "package",
	KeyVoid:         "void",
	KeyType:         "type",
	KeyMap:          "map",
	OpBitOr:         "|",
	OpBitAnd:        "&",
	OpPlus:          "+",
	OpSub:           "=",
}

var keyMap = map[string]TokenType{
	"byte":     KeyByte,
	"sbyte":    KeySByte,
	"int16":    KeyInt16,
	"uint16":   KeyUInt16,
	"int32":    KeyInt32,
	"uint32":   KeyUInt32,
	"int64":    KeyInt64,
	"uint64":   KeyUInt64,
	"float32":  KeyFloat32,
	"float64":  KeyFloat64,
	"string":   KeyString,
	"bool":     KeyBool,
	"enum":     KeyEnum,
	"struct":   KeyStruct,
	"table":    KeyTable,
	"contract": KeyContract,
	"using":    KeyImport,
	"package":  KeyPackage,
	"void":     KeyVoid,
	"throws":   KeyThrows,
	"type":     KeyType,
	"map":      KeyMap,
}

//String implement fmt.Stringer interface
func (token TokenType) String() string {
	if token > 0 {
		return string(token)
	}

	return tokenName[token]
}

//Position position of source code file
type Position struct {
	FileName string //script file name
	Lines    int    //line number, starting at 1
	Column   int    //column number, starting at 1 (character count per line)
}

//ShortName get the source code file short name
func (pos Position) ShortName() string {
	return filepath.Base(pos.FileName)
}

func (pos Position) String() string {
	return fmt.Sprintf("%s(%d,%d)", pos.FileName, pos.Lines, pos.Column)
}

//Valid check if the position object is valid
func (pos Position) Valid() bool {
	if pos.Lines != 0 {
		return true
	}

	return false
}

// Token .
type Token struct {
	Type  TokenType   // token type
	Value interface{} //token value
	Start Position    // star position
	End   Position    // end position
}

//_NewToken create new token with type and token value
func _NewToken(t rune, val interface{}) *Token {
	return &Token{
		Type:  TokenType(t),
		Value: val,
	}
}

func (token *Token) String() string {
	if token.Value != nil {
		return fmt.Sprintf(
			"token {{ %s }}\n\tval :%v\n\tstart :%s\n\tend :%s",
			token.Type,
			token.Value,
			token.Start,
			token.End,
		)
	}

	return fmt.Sprintf(
		"token {{ %s }}\n\tstart :%s\n\tend :%s",
		token.Type,
		token.Start,
		token.End,
	)

}

//Lexer gslang lexer implements reading Of Unicode characters
//and tokens from an io.Reader.
type Lexer struct {
	gslogger.Log                   //Mixin log interface
	reader       *bufio.Reader     //input reader
	position     Position          //curr coursor position
	token        *Token            //curr parsed token
	buff         [utf8.UTFMax]byte //buffer length
	buffPos      int               //buff write position
	offset       int               //reader stream offset by byte
	ws           uint64            //ws flags
	curr         rune              //curr utf8 characters
}

//NewLexer create new gslang lexer
func NewLexer(tag string, reader io.Reader) *Lexer {
	return &Lexer{
		Log:    gslogger.Get("gslang[lexer]"),
		reader: bufio.NewReader(reader),
		ws:     1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' ',
		position: Position{
			FileName: tag,
			Lines:    1,
			Column:   1,
		},
		curr: rune(TokenEOF),
	}
}

func (lexer *Lexer) String() string {
	return lexer.position.FileName
}

func (lexer *Lexer) newerror(fmtstring string, args ...interface{}) error {
	return gserrors.Newf(ErrLexer, "[lexer] %s\n\t%s", fmt.Sprintf(fmtstring, args...), lexer.position)
}

//nextChar read next utf-8 character
func (lexer *Lexer) nextChar() error {

	c, err := lexer.reader.ReadByte()

	if err != nil {
		if err == io.EOF {
			lexer.curr = rune(TokenEOF)
			return nil
		}
		return err
	}

	lexer.offset++
	//not ASCII
	if c >= utf8.RuneSelf {
		lexer.buff[0] = c
		lexer.buffPos = 1
		for !utf8.FullRune(lexer.buff[0:lexer.buffPos]) {
			//continue read rest utf8 char bytes
			c, err = lexer.reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					lexer.curr = rune(TokenEOF)
					return nil
				}
				return err
			}

			lexer.buff[lexer.buffPos] = c
			lexer.buffPos++

			gserrors.Assert(
				lexer.buffPos < len(lexer.buff),
				"utf8.UTFMax must << len(lexer.buff)",
			)
		}

		c, width := utf8.DecodeRune(lexer.buff[0:lexer.buffPos])

		if c == utf8.RuneError && width == 1 {
			return lexer.newerror("illegal utf8 character")

		}

		lexer.curr = c
	} else {
		lexer.curr = rune(c)
	}

	lexer.position.Column++

	return nil
}

func isDecimal(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

//next read next gslang token
func (lexer *Lexer) next() (token *Token, err error) {

	if rune(TokenEOF) == lexer.curr {
		if err = lexer.nextChar(); err != nil {
			return
		}
	}

	//skip white space
	for lexer.ws&(1<<uint(lexer.curr)) != 0 {
		//increase lines count
		if lexer.curr == '\n' {
			lexer.position.Column = 1
			lexer.position.Lines++
		}

		if err = lexer.nextChar(); err != nil {
			return
		}
	}

	//check if arrived end of file
	if lexer.curr == rune(TokenEOF) {
		token = _NewToken(rune(TokenEOF), nil)
		token.Start = lexer.position
		token.End = lexer.position
		return
	}

	position := lexer.position

	position.Column = position.Column - 1

	switch {
	case unicode.IsLetter(lexer.curr) || lexer.curr == '_': //scan id
		token, err = lexer.scanID()

		if err == nil {

			id := token.Value.(string)
			//test if the token is true/false constant or key word
			if id == "true" {
				token.Type = TokenTrue
			} else if id == "false" {
				token.Type = TokenFalse
			} else if key, ok := keyMap[id]; ok {
				token.Type = key
			} else {
				if lexer.curr == ':' {
					token.Type = TokenLABEL
					lexer.nextChar()
				}
			}
		}
	case isDecimal(lexer.curr): //scan number [0-9]+([+-](E|e)[0-9]+)

		token, err = lexer.scanNum()
	case '"' == lexer.curr:
		token, err = lexer.scanString('"')
	case '\'' == lexer.curr:
		token, err = lexer.scanString('\'')
	case '/' == lexer.curr:
		err = lexer.nextChar()

		if err == nil {
			//scan comment
			if lexer.curr == '/' || lexer.curr == '*' {
				token, err = lexer.scanComment(lexer.curr)
			} else {
				token = _NewToken(lexer.curr, nil)
			}
		}

	case '-' == lexer.curr:

		err = lexer.nextChar()

		if err == nil {
			//scan comment
			if lexer.curr == '>' {
				token = _NewToken(rune(TokenArrowRight), nil)
				err = lexer.nextChar()
			} else {
				token = _NewToken(rune(OpSub), nil)
			}
		}
	case '+' == lexer.curr:
		err = lexer.nextChar()
		token = _NewToken(rune(OpPlus), nil)
	case '|' == lexer.curr:
		err = lexer.nextChar()
		token = _NewToken(rune(OpBitOr), nil)
	case '&' == lexer.curr:
		err = lexer.nextChar()
		token = _NewToken(rune(OpBitAnd), nil)
	default:
		token = _NewToken(lexer.curr, nil)
		lexer.curr = rune(TokenEOF)
	}

	if err == nil {
		token.Start = position

		position = lexer.position

		position.Column = position.Column - 1

		token.End = position

	}

	return
}

func (lexer *Lexer) scanComment(ch rune) (*Token, error) {
	var buff bytes.Buffer
	// ch == '/' || ch == '*'
	if ch == '/' {
		// line comment
		err := lexer.nextChar() // read character after "//"
		if err != nil {
			return nil, err
		}
		for lexer.curr != '\n' && lexer.curr >= 0 {
			buff.WriteRune(lexer.curr)
			err = lexer.nextChar() // read character after "//"
			if err != nil {
				return nil, err
			}
		}

		return _NewToken(rune(TokenCOMMENT), buff.String()), nil
	}

	err := lexer.nextChar() // read character after "//"
	if err != nil {
		return nil, err
	}

	for {
		if lexer.curr < 0 {
			return nil, lexer.newerror("comment not terminated")
		}

		if lexer.curr == '\n' {
			lexer.position.Column = 0
			lexer.position.Lines++
		}

		ch0 := lexer.curr
		err = lexer.nextChar()
		if err != nil {
			return nil, err
		}
		if ch0 == '*' && lexer.curr == '/' {
			err = lexer.nextChar()
			if err != nil {
				return nil, err
			}
			break
		}

		buff.WriteRune(ch0)
	}

	return _NewToken(rune(TokenCOMMENT), buff.String()), nil

}

func (lexer *Lexer) scanEscape(buff *bytes.Buffer, quote rune) (err error) {
	err = lexer.nextChar() // read character after '/'
	if err != nil {
		return
	}

	switch lexer.curr {
	case quote:
		buff.WriteRune(lexer.curr)
		err = lexer.nextChar()

		if err != nil {
			return
		}
	default:
		err = lexer.newerror("illegal char escape")
	}

	return
}

func (lexer *Lexer) scanString(quote rune) (token *Token, err error) {
	var buff bytes.Buffer
	err = lexer.nextChar() // read character after quote
	if err != nil {
		return nil, err
	}
	for lexer.curr != quote {
		if lexer.curr == '\n' || lexer.curr < 0 {
			err = lexer.newerror("literal not terminated")
			return
		}
		if lexer.curr == '\\' {
			lexer.scanEscape(&buff, quote)
		} else {
			buff.WriteRune(lexer.curr)
			err = lexer.nextChar()
			if err != nil {
				return nil, err
			}
		}
	}

	err = lexer.nextChar()

	if err != nil {
		return nil, err
	}
	token = _NewToken(rune(TokenSTRING), buff.String())

	return
}

//scanID read id token
func (lexer *Lexer) scanID() (token *Token, err error) {
	var buff bytes.Buffer
	for lexer.curr == '_' || unicode.IsLetter(lexer.curr) || unicode.IsDigit(lexer.curr) {
		buff.WriteRune(lexer.curr) //append to buff
		if err = lexer.nextChar(); err != nil {
			return nil, err
		}
	}

	token = _NewToken(rune(TokenID), string(buff.Bytes()))

	return
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}

//scanNum read number token
func (lexer *Lexer) scanNum() (*Token, error) {

	var buff bytes.Buffer

	if lexer.curr == '0' {

		buff.WriteRune(lexer.curr)

		lexer.nextChar()

		if lexer.curr == 'x' || lexer.curr == 'X' {
			buff.WriteRune(lexer.curr)
			lexer.nextChar()
			for digitVal(lexer.curr) < 16 {
				buff.WriteRune(lexer.curr)
				lexer.nextChar()
			}
			if buff.Len() < 3 {
				return nil, lexer.newerror("illegal hexadecimal number")
			}

			val, err := strconv.ParseInt(buff.String(), 0, 64)

			if err != nil {
				return nil, lexer.newerror(err.Error())
			}

			return _NewToken(rune(TokenINT), val), nil
		}
	}

	lexer.scanMantissa(&buff)

	switch lexer.curr {
	case '.', 'e', 'E':

		lexer.scanFraction(&buff)

		lexer.scanExponent(&buff)

		val, err := strconv.ParseFloat(buff.String(), 64)

		if err != nil {
			return nil, lexer.newerror(err.Error())
		}

		return _NewToken(rune(TokenFLOAT), val), nil
	}

	val, err := strconv.ParseInt(buff.String(), 0, 64)

	if err != nil {
		return nil, lexer.newerror(err.Error())
	}
	return _NewToken(rune(TokenINT), val), nil
}

func (lexer *Lexer) scanMantissa(buff *bytes.Buffer) {
	for isDecimal(lexer.curr) {
		buff.WriteRune(lexer.curr)
		lexer.nextChar()
	}
}

func (lexer *Lexer) scanFraction(buff *bytes.Buffer) {
	if lexer.curr == '.' {
		buff.WriteRune(lexer.curr)
		lexer.nextChar()
		lexer.scanMantissa(buff)
	}
}

func (lexer *Lexer) scanExponent(buff *bytes.Buffer) {
	if lexer.curr == 'e' || lexer.curr == 'E' {
		buff.WriteRune(lexer.curr)
		lexer.nextChar()
		if lexer.curr == '-' || lexer.curr == '+' {
			buff.WriteRune(lexer.curr)
			lexer.nextChar()
		}
		lexer.scanMantissa(buff)
	}
}

//Peek peek next token
func (lexer *Lexer) Peek() (token *Token, err error) {
	if lexer.token != nil {

		token = lexer.token
		return
	}

	token, err = lexer.next()
	if err == nil {
		lexer.token = token
	}

	return
}

//Next get next token and move lexer's cursor
func (lexer *Lexer) Next() (token *Token, err error) {

	if lexer.token != nil {
		token, lexer.token = lexer.token, nil
		return
	}

	token, err = lexer.next()

	return
}
