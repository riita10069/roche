package gen_scaffold

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/riita10069/roche/pkg/rochectl/ast"
	"github.com/riita10069/roche/pkg/util/slice"
)

type InterfaceSpec struct {
	Name       string
	ImportPath string
}

var (
	//go:embed template/wire.go.tmpl
	wireTemplate string
	//go:embed template/provider.go.tmpl
	providerTemplate string
)

// GenerateWireContent generate the content of wire.go
func GenerateWireFile(wireDir string, importPath string) (string, error) {
	importList := []string{}
	providerSetList := []string{}

	filePath := wireDir + "/" + "wire.go"
	if _, err := os.Stat(filePath); err == nil {
		// ファイルがすでに存在してたら
		importList, err = ast.FindImportPath(filePath)
		if err != nil {
			return "", err
		}

		// "github.com/google/wire"を無視
		importList = importList[1:]
	}

	// importPathの重複確認
	if ok, _ := slice.Contains(importPath, importList); !ok {
		importList = append(importList, importPath)
	}

	// importPathからproviderSetListを作成
	for _, path := range importList {
		paths := strings.Split(path, "/")
		pkgName := paths[len(paths)-1]

		providerSetList = append(providerSetList, pkgName+"."+"Set")
	}

	tpl, err := template.New("wireTemplate").Parse(wireTemplate)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"importList":      importList,
		"providerSetList": providerSetList,
	}

	var buf bytes.Buffer

	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GenerateProviderContent generate the content of provider.go
func GenerateProviderFile(providerDir string, providerName string, bindMap map[string]*ast.InterfaceSpec) (string, error) {
	importList := []string{}
	providerList := []string{}
	bMap := map[string]*ast.InterfaceSpec{}

	filePath := providerDir + "/" + "provider.go"
	if _, err := os.Stat(filePath); err == nil {
		// ファイルがすでに存在してたら
		importList, err = ast.FindImportPath(filePath)
		if err != nil {
			return "", err
		}

		// "github.com/google/wire"を無視
		importList = importList[1:]

		providerList, err = ast.FindProviderName(filePath)
		if err != nil {
			return "", err
		}

		bMap, err = ast.FindWireBind(filePath)
		if err != nil {
			return "", err
		}
	}

	if ok, _ := slice.Contains(providerName, providerList); !ok {
		providerList = append(providerList, providerName)
	}

	// bindMapを再作成
	for structName, interfaceSpec := range bindMap {
		bMap[structName] = interfaceSpec
	}

	// bMapからimportListにパスを追加
	for _, interfaceSpec := range bMap {
		importPath := interfaceSpec.ImportPath
		if ok, _ := slice.Contains(importPath, importList); !ok {
			importList = append(importList, importPath)
		}
	}

	// bMapをwire.Bind(new(interface), new(*struct)) の形式にする
	bindList := []string{}
	for structName, interfaceSpec := range bMap {
		paths := strings.Split(interfaceSpec.ImportPath, "/")
		importName := paths[len(paths)-1]

		interfaceName := importName + "." + interfaceSpec.Name
		bind := fmt.Sprintf("wire.Bind(new(%s), new(*%s))", interfaceName, structName)
		bindList = append(bindList, bind)
	}

	tpl, err := template.New("providerTemplate").Parse(providerTemplate)
	if err != nil {
		return "", err
	}

	paths := strings.Split(providerDir, "/")
	pkgName := paths[len(paths)-1]

	data := map[string]interface{}{
		"pkgName":      pkgName,
		"importList":   importList,
		"providerList": providerList,
		"bindList":     bindList,
	}

	var buf bytes.Buffer

	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
