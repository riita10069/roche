package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// FindImportPath find importPath list
func FindImportPath(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.Mode(0))
	if err != nil {
		return nil, err
	}

	importList := []string{}
	for _, importSpec := range f.Imports {
		importPath := strings.Replace(importSpec.Path.Value, "\"", "", -1)
		importList = append(importList, importPath)
	}

	return importList, nil
}

// FindProviderName find provider name list
func FindProviderName(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.Mode(0))
	if err != nil {
		return nil, err
	}

	var genDecl *ast.GenDecl
	for _, decl := range f.Decls {
		d, ok := decl.(*ast.GenDecl)
		if ok && d.Tok == token.VAR {
			genDecl = d
		}
	}

	if genDecl == nil {
		return nil, err
	}

	var valueSpec *ast.ValueSpec
	for _, spec := range genDecl.Specs {
		s, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for _, name := range s.Names {
			if name.Name != "Set" {
				continue
			}
		}

		valueSpec = s
	}

	if valueSpec == nil {
		return nil, err
	}

	var callExpr *ast.CallExpr
	for _, value := range valueSpec.Values {
		expr, ok := value.(*ast.CallExpr)
		if !ok {
			continue
		}

		fun, ok := expr.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		xIdent, ok := fun.X.(*ast.Ident)
		if !ok {
			continue
		}

		if xIdent.Name != "wire" || fun.Sel.Name != "NewSet" {
			continue
		}

		callExpr = expr
	}

	if callExpr == nil {
		return nil, err
	}

	providerNameList := []string{}
	for _, arg := range callExpr.Args {
		ident, ok := arg.(*ast.Ident)
		if !ok {
			continue
		}
		providerNameList = append(providerNameList, ident.Name)
	}

	return providerNameList, nil
}

type InterfaceSpec struct {
	Name       string
	ImportPath string
}

// FindWireBind find elements of wire.Bind()
func FindWireBind(filePath string) (map[string]*InterfaceSpec, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.Mode(0))
	if err != nil {
		return nil, err
	}

	var genDecl *ast.GenDecl
	for _, decl := range f.Decls {
		d, ok := decl.(*ast.GenDecl)
		if ok && d.Tok == token.VAR {
			genDecl = d
		}
	}

	if genDecl == nil {
		return nil, err
	}

	var valueSpec *ast.ValueSpec
	for _, spec := range genDecl.Specs {
		s, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for _, name := range s.Names {
			if name.Name != "Set" {
				continue
			}
		}

		valueSpec = s
	}

	if valueSpec == nil {
		return nil, err
	}

	var callExpr *ast.CallExpr
	for _, value := range valueSpec.Values {
		expr, ok := value.(*ast.CallExpr)
		if !ok {
			continue
		}

		fun, ok := expr.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		xIdent, ok := fun.X.(*ast.Ident)
		if !ok {
			continue
		}

		if xIdent.Name != "wire" || fun.Sel.Name != "NewSet" {
			continue
		}

		callExpr = expr
	}

	if callExpr == nil {
		return nil, err
	}

	// wire.Bindとなるast.SelectorExprを取得
	var bindCallExpr []*ast.CallExpr
	for _, arg := range callExpr.Args {
		cExpr, ok := arg.(*ast.CallExpr)
		if !ok {
			continue
		}

		bindCallExpr = append(bindCallExpr, cExpr)
	}

	// argsを取得 1つめがstruct、2つ目がinterface用
	importMap := map[string]string{}
	for _, importSpec := range f.Imports {
		importPath := strings.Replace(importSpec.Path.Value, "\"", "", -1)

		paths := strings.Split(importPath, "/")
		importName := paths[len(paths)-1]

		importMap[importName] = importPath
	}

	bindMap := map[string]*InterfaceSpec{}
	for _, cExpr := range bindCallExpr {
		// interface
		interfaceCallExpr, _ := cExpr.Args[0].(*ast.CallExpr)
		interfaceSelectorExpr, _ := interfaceCallExpr.Args[0].(*ast.SelectorExpr)
		importNameIdent, _ := interfaceSelectorExpr.X.(*ast.Ident)

		importName := importNameIdent.Name
		interfaceName := interfaceSelectorExpr.Sel.Name

		interfaceSpec := InterfaceSpec{
			Name:       interfaceName,
			ImportPath: importMap[importName],
		}

		// struct
		structCallExpr, _ := cExpr.Args[1].(*ast.CallExpr)
		starExpr, _ := structCallExpr.Args[0].(*ast.StarExpr)

		structIdent, _ := starExpr.X.(*ast.Ident)
		structName := structIdent.Name

		bindMap[structName] = &interfaceSpec
	}

	return bindMap, nil
}
