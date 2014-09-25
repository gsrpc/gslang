package ast

import "testing"

func TestASTCreate(t *testing.T) {
	pkg := NewPackage("test")

	_, err := pkg.NewScript("test.gs")

	if err != nil {
		t.Fatal(err)
	}

	_, err = pkg.NewScript("test.gs")

	if err == nil {
		t.Fatal("create duplicate script test fault")
	}
}
