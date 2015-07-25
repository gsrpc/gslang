package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslang"
	"github.com/gsdocker/gslogger"
)

func TestToken(t *testing.T) {

	content, err := ioutil.ReadFile("test.gs")

	if err != nil {
		t.Fatal(err)
	}

	lexer := gslang.NewLexer("mem", bytes.NewBuffer(content))

	for {

		token, err := lexer.Next()

		if err != nil {
			t.Fatal(err)
		}

		if token.Type == gslang.TokenEOF {
			break
		}

		fmt.Printf("token %s\n", token)
	}
}

func TestParser(t *testing.T) {

	defer gslogger.Join()

	compiler := gslang.NewCompiler()

	err := compiler.Compile("test.gs", gslang.HandleParseError(func(err error, position gslang.Position, msg string) {
		gserrors.Panicf(err, "parse %s error\n\t%s", position, msg)
	}))

	if err != nil {
		t.Fatal(err)
	}
}
