package gen_scaffold

import (
	"github.com/dave/jennifer/jen"
	"github.com/riita10069/roche/pkg/util/slice"
	"go/ast"
)

func GenerateModel(name string, structAst *ast.StructType, marker string) *jen.File {
	f := jen.NewFile("model")

	f.Commentf("+%s", marker)

	// create fields of struct
	var codes []jen.Code
	var code *jen.Statement
	for _, field := range structAst.Fields.List {
		for _, nameIdent := range field.Names {
			// 要素の型名
			switch field.Type.(type) {
			// 別パッケージの型を利用している場合
			case *ast.SelectorExpr:
				code = jen.Id(nameIdent.Name)
				selectorExpr, _ := field.Type.(*ast.SelectorExpr)
				xIdent, _ := selectorExpr.X.(*ast.Ident)
				if xIdent.Name == "protoimpl" {
					continue
				}
				code.Id(xIdent.Name + "." + selectorExpr.Sel.Name)
			// 組み込みまたはどうパッケージ内の型を利用している場合
			case *ast.Ident:
				ident, _ := field.Type.(*ast.Ident)
				code = jen.Id(nameIdent.Name)
				code.Id(ident.Name)
			case *ast.StarExpr:
				code = jen.Id(nameIdent.Name + "ID")
				code.Int64()
			}
			codes = append(codes, code)
		}
	}

	// create struct
	f.Type().Id(name).Struct(codes...)

	return f
}

func isBuiltInType(thisType string) (bool, error) {
	var builtArray []string = []string{"bool", "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "float32", "float64", "complex64", "complex128", "byte", "uint", "rune", "uintptr", "error", "string"}
	return slice.Contains(thisType, builtArray)
}
