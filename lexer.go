package gslang

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslogger"
)

//Lexer errors
var (
	ErrLexer = errors.New("lexer error")
)

//Token types
const (
	TokenEOF rune = -(iota + 1)
	TokenID
	TokenINT
	TokenFLOAT
	TokenTrue
	TokenFalse
	TokenSTRING
	TokenCOMMENT
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
	KeyBool
	KeyEnum
	KeyStruct
	KeyTable
	KeyContract
)

var tokenName = map[rune]string{
	TokenEOF:     "EOF",
	TokenID:      "ID",
	TokenINT:     "INT",
	TokenFLOAT:   "FLOAT",
	TokenSTRING:  "STRING",
	TokenCOMMENT: "COMMENT",
	KeyByte:      "byte",
	KeySByte:     "sbyte",
	KeyInt16:     "int16",
	KeyUInt16:    "uint16",
	KeyInt32:     "int32",
	KeyUInt32:    "uint32",
	KeyInt64:     "int64",
	KeyUInt64:    "uint64",
	KeyFloat32:   "float32",
	KeyFloat64:   "float64",
	KeyBool:      "bool",
	KeyEnum:      "enum",
	KeyStruct:    "struct",
	KeyTable:     "table",
	KeyContract:  "contract",
}

var keyMap = map[string]rune{
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
	"bool":     KeyBool,
	"enum":     KeyEnum,
	"struct":   KeyStruct,
	"table":    KeyTable,
	"contract": KeyContract,
}

//TokenName get token string
func TokenName(token rune) string {
	if token > 0 {
		return string(token)
	}

	return tokenName[token]
}

//Token the gslang token object
type Token struct {
	Type     rune        //token type
	Value    interface{} //token value
	Position Position    //token position in source file
}

//NewToken create new token with type and token value
func NewToken(t rune, val interface{}) *Token {
	return &Token{
		Type:  t,
		Value: val,
	}
}

func (token *Token) String() string {
	if token.Value != nil {
		return fmt.Sprintf(
			"token[%s]\n\tval :%v\n\tpos :%s",
			TokenName(token.Type),
			token.Value,
			token.Position,
		)
	}

	return fmt.Sprintf(
		"token[%s]\n\tpos :%s",
		TokenName(token.Type),
		token.Position,
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
func NewLexer(filename string, reader io.Reader) *Lexer {
	return &Lexer{
		Log:    gslogger.Get("gslang[lexer]"),
		reader: bufio.NewReader(reader),
		ws:     1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' ',
		position: Position{
			FileName: filename,
			Lines:    1,
			Column:   1,
		},
		curr: TokenEOF,
	}
}

func (lexer *Lexer) newerror(fmtstring string, args ...interface{}) error {
	return gserrors.Newf(ErrLexer, "[lexer] %s\n\t%s", fmt.Sprintf(fmtstring, args...), lexer.position)
}

//nextChar read next utf-8 character
func (lexer *Lexer) nextChar() error {
	c, err := lexer.reader.ReadByte()

	if err != nil {
		if err == io.EOF {
			lexer.curr = TokenEOF
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
					lexer.curr = TokenEOF
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

	position := lexer.position

	if TokenEOF == lexer.curr {
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
	if lexer.curr == TokenEOF {
		token = NewToken(TokenEOF, nil)
		token.Position = position
		return
	}
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
			} else {
				if key, ok := keyMap[id]; ok {

					token.Type = key
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
				token = NewToken(lexer.curr, nil)
			}
		}

	default:
		token = NewToken(lexer.curr, nil)
		lexer.curr = TokenEOF
	}

	if err == nil {
		token.Position = position
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

		return NewToken(TokenCOMMENT, buff.String()), nil
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
			lexer.position.Column = 1
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

	return NewToken(TokenCOMMENT, buff.String()), nil

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

	token = NewToken(TokenSTRING, buff.String())

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

	token = NewToken(TokenID, string(buff.Bytes()))

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

	lexer.D("scan num")

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
			lexer.D("scan num finish :%s", buff.String())
			return NewToken(TokenINT, val), nil
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

		lexer.D("scan num finish :%s", buff.String())

		return NewToken(TokenFLOAT, val), nil
	}

	val, err := strconv.ParseInt(buff.String(), 0, 64)

	if err != nil {
		return nil, lexer.newerror(err.Error())
	}
	lexer.D("scan num finish :%s", buff.String())
	return NewToken(TokenINT, val), nil
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

// //Peek peek next token
// func (lexer *Lexer) Peek() (token *Token, err error) {
// 	if lexer.token != nil {
// 		token = lexer.token
// 		return
// 	}
//
// 	token, err = lexer.next()
//
// 	if err != nil {
// 		lexer.token = token
// 	}
//
// 	return
// }
//
// //Next get next token and move lexer's cursor
// func (lexer *Lexer) Next() (token *Token, err error) {
//
// 	if lexer.token != nil {
// 		token, lexer.token = lexer.token, nil
// 		return
// 	}
//
// 	token, err = lexer.next()
//
// 	return
// }
