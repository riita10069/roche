package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func FindInterface(interfaceName string, filePath string) *ast.InterfaceType {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.Mode(0))
	if err != nil {
		return nil
	}

	interfaceASTMap := map[string]*ast.InterfaceType{}
	getInterfaceAstHash(f, interfaceASTMap)

	target := interfaceASTMap[interfaceName]
	return target
}

func getInterfaceAstHash(f *ast.File, interfaceASTMap map[string]*ast.InterfaceType) {
	for _, decl := range f.Decls {
		// GenDeclであるかを判定
		d, ok := decl.(*ast.GenDecl)
		if !ok || d.Tok != token.TYPE {
			continue
		}

		for _, spec := range d.Specs {
			s, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			t, ok := s.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			interfaceASTMap[s.Name.Name] = t
		}
	}
}

