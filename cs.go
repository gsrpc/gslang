package gslang

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslogger"
)

//ErrComileS public errors
var (
	ErrCompileS = errors.New("CompileS error")
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
		gserrors.Newf(ErrCompileS, "must set GOPATH first")
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
		gserrors.Newf(ErrCompileS, stream.String())
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

func (cs *CompileS) errorf(position Position, fmtstring string, args ...interface{}) {
	gserrors.Panicf(
		ErrParse,
		fmt.Sprintf(
			"parse %s error : %s",
			position,
			fmt.Sprintf(fmtstring, args...),
		),
	)
}

//Compile 编译指定的包
func (cs *CompileS) Compile(packageName string) (pkg *ast.Package, err error) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(gserrors.GSError); ok {
				err = e.(error)
			} else {
				err = gserrors.New(e.(error))
			}
		}
	}()

	defer gserrors.Ensure(func() bool {
		if err == nil {
			return pkg != nil
		}

		return true

	}, "if err == nil ,the return pkg param can't be nil")

	//跳过已经加载的包
	if loaded, ok := cs.Loaded[packageName]; ok {
		//cs.D("skip compiled package :%s ", packageName)
		pkg = loaded
		return
	}
	//循环引用检测
	cs.circularRefCheck(packageName)

	fullPath := cs.searchPackage(packageName)

	pkg = ast.NewPackage(packageName)

	cs.loading = append(cs.loading, pkg)

	err = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() && path != fullPath {
			return filepath.SkipDir
		}

		if filepath.Ext(info.Name()) != ".gs" {
			return nil
		}

		script, err := cs.parse(pkg, path)

		if err == nil {
			setFilePath(script, path)
		}

		return err
	})

	if err != nil {
		cs.loading = cs.loading[:len(cs.loading)-1]
		return
	}

	cs.link(pkg)
	//
	// cs._MoveAttr(pkg)

	cs.loading = cs.loading[:len(cs.loading)-1]
	cs.Loaded[packageName] = pkg
	return
}
