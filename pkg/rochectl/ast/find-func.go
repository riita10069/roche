package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func FindFunc(filepath string) map[string]*ast.FuncDecl {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filepath, nil, parser.Mode(0))
	if err != nil {
		return nil
	}
	funcAstHash := getFuncAstHash(f)

	return funcAstHash
}

func getFuncAstHash(f *ast.File) map[string]*ast.FuncDecl {
	ret := map[string]*ast.FuncDecl{}
	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok  {
			continue
		}
		ret[funcDecl.Name.Name] = funcDecl
	}

	return ret
}