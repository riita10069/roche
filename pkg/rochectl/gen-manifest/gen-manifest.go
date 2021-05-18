package gen_manifest

import (
	"errors"
	"github.com/riita10069/roche/pkg/rochectl/config"
	"github.com/riita10069/roche/pkg/rochectl/file"
	"github.com/riita10069/roche/pkg/rochectl/gen-manifest/tmpl"
	"strings"
	"text/template"
)

func GenerateDeployment(name string, cnf *config.Config) error {
	parse, err := template.New("deployment.yml").Parse(tmpl.DeploymentTemplate)
	if err != nil {
		return errors.New("cannot parse deployment template")
	}
	hash := map[string]interface{}{
		"Name":     name,
		"Registry": cnf.ImageRegistry,
	}

	writer := new(strings.Builder)
	if err := parse.Execute(writer, hash); err != nil {
		errors.New("template parse error")
	}

	err = file.CreateAndWrite(writer.String(), cnf.ManifestsDir+"/"+name+"/base/deployment.yaml")
	if err != nil {
		return err
	}
	return nil
}

func GenerateHpa(name string, cnf *config.Config) error {
	parse, err := template.New("hpa.yml").Parse(tmpl.HpaTemplate)
	if err != nil {
		return errors.New("cannot parse hpa template")
	}
	hash := map[string]interface{}{
		"Name":     name,
	}

	writer := new(strings.Builder)
	if err := parse.Execute(writer, hash); err != nil {
		errors.New("template parse error")
	}

	err = file.CreateAndWrite(writer.String(), cnf.ManifestsDir + "/" + name + "/base/hpa.yaml")
	if err != nil {
		return err
	}
	return nil
}

func GenerateNodePort(name string, cnf *config.Config) error {
	parse, err := template.New("service.yml").Parse(tmpl.NodePortTemplate)
	if err != nil {
		return errors.New("cannot parse node port template")
	}
	hash := map[string]interface{}{
		"Name":     name,
	}

	writer := new(strings.Builder)
	if err := parse.Execute(writer, hash); err != nil {
		errors.New("template parse error")
	}

	err = file.CreateAndWrite(writer.String(), cnf.ManifestsDir+"/"+name+"/base/service.yaml")
	if err != nil {
		return err
	}
	return nil
}
