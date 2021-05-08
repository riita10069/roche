package ast

import (
	"fmt"
	"go/ast"
	"go/token"
)

func PrintFile(f ast.File)  {
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
