package config

import (
	"github.com/riita10069/roche/pkg/util"
)

type Config struct {
	ModuleName    string
	ServerDir     string
	DiDir         string
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

func (c Config) GetInfraModelFilePath(name string) string {
	return c.InfraModelDir + "/" + util.CamelToSnake(name) + ".go"
}

func (c Config) GetInfraRepoFilePath(name string) string {
	return c.RepoDir + "/" + util.CamelToSnake(name) + ".go"
}

func (c Config) GetServerFilePath(name string) string {
	return c.ServerDir + "/" + util.CamelToSnake(name) + ".go"
}
