package roche_gen_toml

import (
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/riita10069/roche/pkg/roche/file"
	"github.com/spf13/cobra"
)

const content =
`ModuleName    = "github.com/repository/package"
ServerDir     = "app/server"
UsecaseDir    = "usecase"
DomainRepoDir = "domain/repository"
EntityDir     = "domain/entity"
InfraModelDir = "infra/model"
RepoDir       = "infra/repository"
ProtoDir      = "api/proto"
PbGoDir       = "api"
ManifestsDir  = "manifests"
ImageRegistry = "example.io/company"
`

func NewTomlCommand(ctx *grapicmd.Ctx) *cobra.Command {
	tomlCmd := &cobra.Command{
		Use:           "toml",
		Short:         "initialize roche.toml file",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return MakeRocheToml()
		},
	}
	return tomlCmd
}

func MakeRocheToml() error {
	err := file.CreateAndWrite(content, "roche.toml")
	if err != nil {
		return err
	}
	return nil
}
