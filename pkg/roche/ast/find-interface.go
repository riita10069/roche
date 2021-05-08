package ast

import (
	"github.com/riita10069/roche/pkg/roche/config"
	"github.com/riita10069/roche/pkg/util"
	"go/ast"
	"go/parser"
	"go/token"
)

func FindInterface(name string, cnf *config.Config) *ast.InterfaceType {
	interfaceName := name + "ServiceServer"
	filePath := cnf.PbGoDir + "/" + util.CamelToSnake(name) + ".pb.go"
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

