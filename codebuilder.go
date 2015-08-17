package gslang

import (
	"fmt"
	"io"
	"text/template"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslang/ast"
	"github.com/gsdocker/gslang/lexer"
)

// CodeBuilder .
type CodeBuilder interface {
	CreateTable(writer io.Writer, table *ast.Table)
	CreateAnnotation(writer io.Writer, table *ast.Table)
	CreateEnum(writer io.Writer, enum *ast.Enum)
	CreateContract(writer io.Writer, contract *ast.Contract)
	CreateUsing(writer io.Writer, text string)
	Reset()
}

// _CodeBuilder .
type _CodeBuilder struct {
	tpl     *template.Template // template
	usings  []string           // usings
	builtin map[lexer.TokenType]string
}

// NewCodeBuilder parse code generate template
func NewCodeBuilder(text string, builtin map[lexer.TokenType]string) (CodeBuilder, error) {

	builder := &_CodeBuilder{
		builtin: builtin,
	}

	funcs := template.FuncMap{
		"enumType": builder.enumType,
	}

	tpl, err := template.New("_CodeBuilder").Funcs(funcs).Parse(text)

	if err != nil {
		return nil, err
	}

	builder.tpl = tpl

	return builder, nil
}

func (builder *_CodeBuilder) enumType(enum *ast.Enum) string {

	_, ok := FindAnnotation(enum, "gslang.Flag")

	if ok {
		return builder.builtin[lexer.KeyUInt32]
	}

	return builder.builtin[lexer.KeyByte]
}

func (builder *_CodeBuilder) Reset() {
	builder.usings = nil
}

// CreateTable .
func (builder *_CodeBuilder) CreateUsing(writer io.Writer, text string) {
	for _, using := range builder.usings {
		if using == text {
			return
		}
	}

	builder.usings = append(builder.usings, text)

	writer.Write([]byte(fmt.Sprintf("%s\n", text)))
}

func (builder *_CodeBuilder) CreateTable(writer io.Writer, table *ast.Table) {
	if err := builder.tpl.ExecuteTemplate(writer, "table", table); err != nil {
		gserrors.Panicf(err, "execute code template(table) error")
	}
}

// CreateAnnotation .
func (builder *_CodeBuilder) CreateAnnotation(writer io.Writer, table *ast.Table) {
	if err := builder.tpl.ExecuteTemplate(writer, "annotation", table); err != nil {
		gserrors.Panicf(err, "execute code template(annotation) error")
	}
}

// CreateEnum .
func (builder *_CodeBuilder) CreateEnum(writer io.Writer, enum *ast.Enum) {
	if err := builder.tpl.ExecuteTemplate(writer, "enum", enum); err != nil {
		gserrors.Panicf(err, "execute code template(enum) error")
	}
}

// CreateContract .
func (builder *_CodeBuilder) CreateContract(writer io.Writer, contract *ast.Contract) {
	if err := builder.tpl.ExecuteTemplate(writer, "contract", contract); err != nil {
		gserrors.Panicf(err, "execute code template(contract) error")
	}
}