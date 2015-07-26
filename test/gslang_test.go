package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslang"
	"github.com/gsdocker/gslang/lexer"
	"github.com/gsdocker/gslogger"
)

func TestToken(t *testing.T) {

	content, err := ioutil.ReadFile("test.gs")

	if err != nil {
		t.Fatal(err)
	}

	tokenizer := lexer.NewLexer("mem", bytes.NewBuffer(content))

	for {

		token, err := tokenizer.Next()

		if err != nil {
			t.Fatal(err)
		}

		if token.Type == lexer.TokenEOF {
			break
		}

		fmt.Printf("token %s\n", token)
	}
}

func TestParser(t *testing.T) {

	defer gslogger.Join()

	compiler := gslang.NewCompiler()

	err := compiler.Compile("test.gs", gslang.HandleParseError(func(err error, position lexer.Position, msg string) {
		gserrors.Panicf(err, "parse %s error\n\t%s", position, msg)
	}))

	if err != nil {
		t.Fatal(err)
	}
}
