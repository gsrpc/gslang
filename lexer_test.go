package gslang

import (
	"bytes"
	"testing"

	"github.com/gsdocker/gserrors"
)

func tokenCheck(lexer *Lexer, expect rune, val interface{}) {
	token, err := lexer.Next()

	if err != nil {
		panic(err)
	}

	switch expect {
	case TokenINT:
		switch val.(type) {
		case int:
			if token.Value.(int64) != int64(val.(int)) {
				gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
			}
		case int32:
			if token.Value.(int64) != int64(val.(int32)) {
				gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
			}
		case int64:
			if token.Value.(int64) != val.(int64) {
				gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
			}
		default:
			gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
		}
	case TokenFLOAT:
		switch val.(type) {
		case float32:
			if token.Value.(float64) != float64(val.(float32)) {
				gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
			}
		case float64:
			if token.Value.(float64) != val.(float64) {
				gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
			}
		default:
			gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
		}
	case TokenID:
		if token.Type != TokenID {
			gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
		}
		switch val.(type) {
		case string:
			if token.Value.(string) != val.(string) {
				gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
			}
		default:
			gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
		}
	case TokenSTRING:
		if token.Type != TokenSTRING {
			gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
		}
		switch val.(type) {
		case string:
			if token.Value.(string) != val.(string) {
				gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
			}
		default:
			gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
		}
	case TokenCOMMENT:
		if token.Type != TokenCOMMENT {
			gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
		}
		switch val.(type) {
		case string:
			if token.Value.(string) != val.(string) {
				gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
			}
		default:
			gserrors.Panicf(nil, "scan token %v err, got \n\t%s", val, token)
		}
	default:
		if token.Type != expect {
			gserrors.Panicf(nil, "scan token %s err, got \n\t%s", TokenName(expect), token)
		}
	}

}

func TestLexerNumber(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal(e.(error))
		}
	}()

	var buff bytes.Buffer

	buff.WriteString(`
		12.5 125
		0x100 0X200

		1.2E+10 12.5E-10 12.5E2
    `)

	lexer := NewLexer("test", &buff)

	tokenCheck(lexer, TokenFLOAT, 12.5)
	tokenCheck(lexer, TokenINT, 125)
	tokenCheck(lexer, TokenINT, 0x100)
	tokenCheck(lexer, TokenINT, 0X200)
	tokenCheck(lexer, TokenFLOAT, 1.2E+10)
	tokenCheck(lexer, TokenFLOAT, 12.5E-10)
	tokenCheck(lexer, TokenFLOAT, 12.5E2)
}

func TestLexerID(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal(e.(error))
		}
	}()

	var buff bytes.Buffer

	buff.WriteString(`
		hello world



		 _hello123
	`)

	lexer := NewLexer("test", &buff)

	tokenCheck(lexer, TokenID, "hello")
	tokenCheck(lexer, TokenID, "world")
	tokenCheck(lexer, TokenID, "_hello123")
}

func TestLexerString(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal(e.(error))
		}
	}()

	var buff bytes.Buffer

	buff.WriteString(`
		"hell \"world\""
	`)

	lexer := NewLexer("test", &buff)

	tokenCheck(lexer, TokenSTRING, "hell \"world\"")
}

func TestLexerComment(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal(e.(error))
		}
	}()

	var buff bytes.Buffer

	buff.WriteString(
		`//"hell "world""
		/*********/
		[]
`,
	)

	lexer := NewLexer("test", &buff)

	tokenCheck(lexer, TokenCOMMENT, `"hell "world""`)
	tokenCheck(lexer, TokenCOMMENT, "*******")

	tokenCheck(lexer, '[', nil)
	tokenCheck(lexer, ']', nil)
}

func TestLexerKeyWorld(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal(e.(error))
		}
	}()

	var buff bytes.Buffer

	buff.WriteString(`
		int32 int64

		int32int64
		`)

	lexer := NewLexer("test", &buff)

	tokenCheck(lexer, KeyInt32, nil)
	tokenCheck(lexer, KeyInt64, nil)
	tokenCheck(lexer, TokenID, "int32int64")
}

func TestLexerPeek(t *testing.T) {

	var buff bytes.Buffer

	buff.WriteString(`
		int32 float

		int32int64
		`)

	lexer := NewLexer("test", &buff)

	token, _ := lexer.Peek()
	token2, _ := lexer.Peek()

	if token.Type != token2.Type {
		t.Fatal("peek test failed")
	}
}
