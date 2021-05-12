package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func PrintFile(filepath string) {
	f, err := parser.ParseFile(token.NewFileSet(), filepath, nil, parser.Mode(0))
	if err != nil {
		fmt.Errorf("cannot parse" + filepath + "    **debug ast.printfile")
	}
	for _, d := range f.Decls {
		ast.Print(token.NewFileSet(), d)
		fmt.Println()
	}
}

func PrintStruct(structAst ast.StructType)  {
	for _, d := range structAst.Fields.List {
		ast.Print(token.NewFileSet(), d)
		fmt.Println()
	}
}
