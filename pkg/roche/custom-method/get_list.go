package custom_method

import (
	goAst "go/ast"
	"strings"
)

func GetList(targetInterface *goAst.InterfaceType, handlerFuncList map[string]*goAst.FuncDecl, name string) []string {
	var funcCode string
	var funcCodeList []string
	for _, signature := range targetInterface.Methods.List {
		_, ok := handlerFuncList[signature.Names[0].Name]
		if ok {
			continue
		}
		paramParen := getFuncParameterCode(signature)
		returnType := getFuncReturnCode(signature)
		funcCode = "func (s *" + name + "ServiceServerImpl) " + signature.Names[0].Name + paramParen + " " + returnType + "{}\n"
		funcCodeList = append(funcCodeList, funcCode)
	}

	return funcCodeList
}

func getFuncParameterCode(signature *goAst.Field) string {
	var parameterCode []string
	for _, param := range signature.Type.(*goAst.FuncType).Params.List {
		switch p := param.Type.(type) {
		case *goAst.SelectorExpr:
			parameterCode = append(parameterCode, p.X.(*goAst.Ident).Name + "." + p.Sel.Name)
		case *goAst.StarExpr:
			parameterCode = append(parameterCode, p.X.(*goAst.Ident).Name)
		case *goAst.Ident:
			parameterCode = append(parameterCode, p.Name)
		}
	}
	return "(" + strings.Join(parameterCode, ", ") + ")"
}

func getFuncReturnCode(signature *goAst.Field) string {
	var returnCode []string
	for _, param := range signature.Type.(*goAst.FuncType).Results.List {
		switch p := param.Type.(type) {
		case *goAst.StarExpr:
			switch p.X.(type) {
			case *goAst.SelectorExpr:
				returnCode = append(returnCode, p.X.(*goAst.SelectorExpr).X.(*goAst.Ident).Name + "." + p.X.(*goAst.SelectorExpr).Sel.Name)
			case *goAst.Ident:
				returnCode = append(returnCode, p.X.(*goAst.Ident).Name)
			}
		case *goAst.Ident:
			returnCode = append(returnCode, p.Name)
		}
	}
	var ret string
	if len(signature.Type.(*goAst.FuncType).Params.List) > 1 {
		ret = "(" + strings.Join(returnCode, ", ") + ")"
	} else {
		ret = strings.Join(returnCode, ", ")
	}
	return ret
}

