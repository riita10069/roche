package roche_gen_toml

import (
	"fmt"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/riita10069/roche/pkg/roche/config"
	"github.com/spf13/cobra"
	"io/ioutil"
)

const data =
`
ModuleName    = "github.com/repository/package"
ServerDir     = "app/server"
UsecaseDir    = "domain/usecase"
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
			err := ioutil.WriteFile("roche.toml", []byte(data), 0664)
			if err != nil {
				fmt.Println(err)
			}
			return nil
		},
	}
	return tomlCmd
}

