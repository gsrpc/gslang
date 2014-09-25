package gslang

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslogger"
)

func setFilePath(node ast.Node, fullPath string) {
	node.NewExtra("FilePath", fullPath)
}

//CompileS gslang compile service
type CompileS struct {
	gslogger.Log                         //Mixin log APIs
	Loaded       map[string]*ast.Package //loaded packages
	loading      []*ast.Package          //loading package path
	goPath       []string                //golang path
}

//NewCompileS create new compile service object
func NewCompileS() *CompileS {
	GOPATH := os.Getenv("GOPATH")

	if GOPATH == "" {
		panic(errors.New("must set GOPATH first"))
	}
	return &CompileS{
		Log:    gslogger.Get("gslang"),
		Loaded: make(map[string]*ast.Package),
		goPath: strings.Split(GOPATH, string(os.PathListSeparator)),
	}
}

func (cs *CompileS) searchPackage(packageName string) string {
	var found []string
	for _, path := range cs.goPath {
		fullpath := filepath.Join(path, "src", packageName)
		fi, err := os.Stat(path)

		if err == nil && fi.IsDir() {
			found = append(found, fullpath)
		}
	}

	if len(found) > 1 {
		var stream bytes.Buffer

		stream.WriteString(fmt.Sprintf("found more than one package named :%s", packageName))

		for i, path := range found {
			stream.WriteString(fmt.Sprintf("\n\t%d) %s", i, path))
		}

		panic(errors.New(stream.String()))
	}

	return found[0]
}

func (cs *CompileS) circularRefCheck(packageName string) {
	var stream bytes.Buffer

	for _, pkg := range cs.loading {
		if pkg.Name() == packageName || stream.Len() != 0 {
			stream.WriteString(fmt.Sprintf("\t%s import\n", pkg.Name()))
		}
	}

	if stream.Len() != 0 {
		panic(fmt.Errorf("circular package import :\n%s\t%s", stream.String(), packageName))
	}
}

//Compile 编译指定的包
func (cs *CompileS) Compile(packageName string) (pkg *ast.Package, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			cs.D("Compile err :%s", err)
		}
	}()
	//跳过已经加载的包
	if _, ok := cs.Loaded[packageName]; ok {
		cs.D("skip compiled package :%s ", packageName)
		return
	}
	//循环引用检测
	cs.circularRefCheck(packageName)

	fullPath := cs.searchPackage(packageName)

	cs.D("loading package :%s\n\tfullpath:%s", packageName, fullPath)

	pkg = ast.NewPackage(packageName)

	err = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() && path != fullPath {
			return filepath.SkipDir
		}

		if filepath.Ext(info.Name()) != ".tap" {
			return nil
		}

		script, err := parse(pkg, path)

		if err == nil {
			setFilePath(script, path)
		}

		return err
	})

	if err != nil {
		return
	}

	cs.loading = append(cs.loading, pkg)

	for _, script := range pkg.Scripts {

		for _, imp := range script.Imports {
			_, err = cs.Compile(imp.Name())

			if err != nil {
				return
			}
		}
	}

	// cs._Link(pkg)
	//
	// cs._MoveAttr(pkg)

	cs.loading = cs.loading[:len(cs.loading)-1]

	cs.Loaded[packageName] = pkg

	return
}
