package gslang

import (
	"fmt"
	"testing"

	"github.com/gsdocker/gslogger"
)

func TestParseImports(t *testing.T) {
	defer gslogger.Join()

	cs := NewCompileS()
	pkg, err := cs.Compile("github.com/gsdocker/gslang/testing/import")

	if err != nil {
		t.Fatal(err)
	}
	script, ok := pkg.Scripts["import.gs"]
	if !ok {
		t.Fatal("loading 'import.gs' err")
	}

	testing, ok := script.Imports["testing"]

	if !ok {
		t.Fatal("test import 'testing' package -- failed")
	}

	comments := Comments(script)

	if comments == nil {
		t.Fatal("test script comments failed")
	}

	for _, comment := range comments {
		fmt.Println("comment", comment.Value.(string))
	}

	comments = Comments(testing)

	if comments == nil {
		t.Fatal("test import comments failed")
	}

	for _, comment := range comments {
		fmt.Println("comment", comment.Value.(string))
	}

}
