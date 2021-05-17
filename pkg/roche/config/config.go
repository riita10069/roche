package config

import "github.com/riita10069/roche/pkg/util"

type Config struct {
	ModuleName    string
	ServerDir     string
	UsecaseDir    string
	DomainRepoDir string
	EntityDir     string
	InfraModelDir string
	RepoDir       string
	ProtoDir      string
	PbGoDir       string
	ManifestsDir  string
	ImageRegistry string
	FindToml      bool
}

func (c Config) GetPbGoFilePath(name string) string {
	return c.PbGoDir + "/" + util.CamelToSnake(name) + ".pb.go"
}

func (c Config) GetEntityFilePath(name string) string {
	return c.EntityDir + "/" + util.CamelToSnake(name) + ".go"
}

func (c Config) GetUsecaseFilePath(name string) string {
	return c.UsecaseDir + "/" + util.CamelToSnake(name) + ".go"
}

func (c Config) GetDomainRepoFilePath(name string) string {
	return c.DomainRepoDir + "/" + util.CamelToSnake(name) + ".go"
}
