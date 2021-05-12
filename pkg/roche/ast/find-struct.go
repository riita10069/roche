package ast

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/riita10069/roche/pkg/util"
	"go/ast"
	"go/parser"
	"go/token"
)

func FindStruct(structName string, filePath string) *ast.StructType {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.Mode(0))
	if err != nil {
		fmt.Println("cannot parse file, filepath is ", filePath)
		fmt.Println(err.Error())
		return nil
	}

	structASTMap := map[string]*ast.StructType{}
	getAstHash(f, structASTMap)
	target := structASTMap[structName]
	return target
}

func getAstHash(f *ast.File, structASTMap map[string]*ast.StructType) {
	for _, decl := range f.Decls {
		// struct型の型定義を表示する

		// GenDeclであるかを判定
		d, ok := decl.(*ast.GenDecl)
		if !ok || d.Tok != token.TYPE {
			continue
		}

		for _, spec := range d.Specs {
			// TypeSpec型か確認
			s, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// StructType型か確認
			t, ok := s.Type.(*ast.StructType)
			if !ok {
				continue
			}

			structASTMap[s.Name.Name] = t
		}
	}
}

func GetPropertyByStructAst(structAst *ast.StructType) ([]string, []string) {
	var property []string
	var propertyType []string

	for _, field := range structAst.Fields.List {
		for _, nameIdent := range field.Names {
			// 要素名


			// 要素の型名
			switch field.Type.(type) {
			// 別パッケージの型を利用している場合
			case *ast.SelectorExpr:
				selectorExpr, _ := field.Type.(*ast.SelectorExpr)
				xIdent, _ := selectorExpr.X.(*ast.Ident)
				if xIdent.Name == "protoimpl" {
					continue
				}
				property = append(property, nameIdent.Name)
				propertyType = append(propertyType, xIdent.Name)

			// 組み込みまたは同パッケージ内の型を利用している場合
			case *ast.Ident:
				ident, _ := field.Type.(*ast.Ident)
				property = append(property, nameIdent.Name)
				propertyType = append(propertyType, ident.Name)
			}
		}
	}
	return property, propertyType
}

func GetPostSignature(property []string, propertyType []string) []jen.Code {
	var postSignature []jen.Code
	for i := range property {
		postSignature = append(postSignature, jen.Id(util.CamelToSnake(property[i])).Id(propertyType[i]))
	}
	return postSignature
}
