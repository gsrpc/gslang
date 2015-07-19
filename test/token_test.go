package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/gsdocker/gslang"
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
