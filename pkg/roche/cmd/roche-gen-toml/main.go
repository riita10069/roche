package roche_gen_toml

import (
	"errors"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/spf13/cobra"
	"io/ioutil"
)

const data =
`
ModuleName    = "github.com/repository/package"
ServerDir     = "app/server"
UsecaseDir    = "usecase"
DomainRepoDir = "domain/repository"
EntityDir     = "domain/entity"
RepoDir       = "infra/repository"
ProtoDir      = "api/proto"
PbGoDir       = "api"
`

func NewTomlCommand(ctx *grapicmd.Ctx) *cobra.Command {
	tomlCmd := &cobra.Command{
		Use:           "toml",
		Short:         "create roche.toml file",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return MakeRocheToml()
		},
	}
	return tomlCmd
}

func MakeRocheToml() error {
	err := ioutil.WriteFile("roche.toml", []byte(data), 0664)
	if err != nil {
		return errors.New("cannot write roche.toml")
	}
	return nil
}
