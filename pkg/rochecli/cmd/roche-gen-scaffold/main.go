package roche_gen_scaffold

import (
	"errors"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/riita10069/roche/pkg/roche/ast"
	"github.com/riita10069/roche/pkg/roche/config"
	"github.com/riita10069/roche/pkg/roche/file"
	gen_scaffold "github.com/riita10069/roche/pkg/roche/gen-scaffold"
	"github.com/riita10069/roche/pkg/util"
	"github.com/spf13/cobra"
)

func NewScaffoldCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	scaffoldCmd := &cobra.Command{
		Use:           "scaffold",
		Short:         "make CRUD code following clean architecture.",
		SilenceErrors: true,
		SilenceUsage:  true,
		Aliases: []string{"s", "sca"},
	}

	scaffoldCmd.AddCommand(NewScaffoldAllCommand(ctx, cnf))
	scaffoldCmd.AddCommand(NewScaffoldModelCommand(ctx, cnf))
	return scaffoldCmd
}

func NewScaffoldAllCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	allCmd := &cobra.Command{
		Use:           "all",
		Short:         "generate all method.",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ctx.IsInsideApp() {
				return errors.New("roche command should be execute inside a roche application directory")
			}
			if !cnf.FindToml {
				return errors.New("For using this roche command, please run following command\nroche toml\nAnd please edit roche.toml According to your project")
			}
			name := args[0]
			pbGoFilePath := cnf.GetPbGoFilePath(name)
			targetStruct := ast.FindStruct(name, pbGoFilePath)
			if targetStruct == nil {
				return errors.New("found "+ pbGoFilePath + " but not found" + name + " struct")
			}
			entityFile := gen_scaffold.GenerateEntity(name, targetStruct)
			file.JenniferToFile(entityFile, cnf.GetEntityFilePath(name))
			usecaseFile, domainRepositoryFile := gen_scaffold.GenerateUsecase(name, targetStruct)
			file.JenniferToFile(usecaseFile, cnf.GetUsecaseFilePath(name))
			file.JenniferToFile(domainRepositoryFile, cnf.GetDomainRepoFilePath(name))
			return nil
		},
	}
	return allCmd
}

func NewScaffoldModelCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	modelCmd := &cobra.Command{
		Use:           "model",
		Short:         "generate model struct.",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			pbGoFilePath := cnf.PbGoDir + "/" + util.CamelToSnake(name) + ".pb.go"
			targetStruct := ast.FindStruct(name, pbGoFilePath)
			if targetStruct == nil {
				return errors.New("found "+ pbGoFilePath + "but not found" + name + " struct")
			}
			err := gen_scaffold.GenerateModel(name, targetStruct, cnf)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return modelCmd
}
