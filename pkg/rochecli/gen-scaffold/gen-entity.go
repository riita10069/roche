package gen_scaffold

import (
	"github.com/dave/jennifer/jen"
	"go/ast"
)


func GenerateEntity(structName string, structAst *ast.StructType) *jen.File {
	f := jen.NewFile("entity")

	// create fields of struct
	var codes []jen.Code

	for _, field := range structAst.Fields.List {
		for _, nameIdent := range field.Names {
			code := jen.Id(nameIdent.Name)
			// 要素の型名
			switch field.Type.(type) {
			// 別パッケージの型を利用している場合
			case *ast.SelectorExpr:
				selectorExpr, _ := field.Type.(*ast.SelectorExpr)
				xIdent, _ := selectorExpr.X.(*ast.Ident)
				if xIdent.Name == "protoimpl" {
					continue
				}
				code.Id(xIdent.Name + "." + selectorExpr.Sel.Name)
			// 組み込みまたはどうパッケージ内の型を利用している場合
			case *ast.Ident:
				ident, _ := field.Type.(*ast.Ident)
				code.Id(ident.Name)
			case *ast.StarExpr:
				code.Id(nameIdent.Name)
			}

			codes = append(codes, code)
		}
	}

	// create struct
	f.Type().Id(structName).Struct(codes...)

	return f
}
