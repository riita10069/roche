package roche_gen_scaffold

import (
	"errors"
	"fmt"
	autoTable "github.com/hourglasshoro/auto-table/pkg"
	autoAutoFile "github.com/hourglasshoro/auto-table/pkg/file"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/naoina/go-stringutil"
	"github.com/riita10069/roche/pkg/rochectl/ast"
	"github.com/riita10069/roche/pkg/rochectl/config"
	"github.com/riita10069/roche/pkg/rochectl/file"
	gen_scaffold "github.com/riita10069/roche/pkg/rochectl/gen-scaffold"
	"github.com/riita10069/roche/pkg/util"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func NewScaffoldCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	scaffoldCmd := &cobra.Command{
		Use:           "scaffold",
		Short:         "make CRUD code following clean architecture.",
		SilenceErrors: true,
		SilenceUsage:  true,
		Aliases:       []string{"s", "sca"},
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
				return errors.New("roche command should be execute inside a rochectl application directory")
			}
			if !cnf.FindToml {
				return errors.New("For using this roche command, please run following command\nrochectl toml\nAnd please edit rochectl.toml According to your project")
			}
			name := util.SnakeToUpperCamel(args[0])
			pfile, err := cmd.PersistentFlags().GetString("pfile")
			if err != nil {
				return err
			}
			if pfile == "default" {
				pfile = name
			}
			pbGoFilePath := cnf.GetPbGoFilePath(pfile)
			targetStruct := ast.FindStruct(name, pbGoFilePath)
			if targetStruct == nil {
				return errors.New("found " + pbGoFilePath + " but not found" + name + " struct")
			}
			entityFile := gen_scaffold.GenerateEntity(name, targetStruct)
			file.JenniferToFile(entityFile, cnf.GetEntityFilePath(name))
			domainRepositoryFile, usecaseFile := gen_scaffold.GenerateUsecase(name, targetStruct)
			file.JenniferToFile(usecaseFile, cnf.GetUsecaseFilePath(name))
			file.JenniferToFile(domainRepositoryFile, cnf.GetDomainRepoFilePath(name))
			infraModelFile := gen_scaffold.GenerateModel(name, targetStruct, ctx.Build.AppName)
			file.JenniferToFile(infraModelFile, cnf.GetInfraModelFilePath(name))
			generator := autoTable.NewGenerator(ctx.Build.AppName)
			files, err := autoAutoFile.GetFiles(&ctx.FS, cnf.InfraModelDir)
			if err != nil {
				return xerrors.Errorf("cannot read infra model files: %w", err)
			}
			sqlMap, err := generator.CreateSQL(files)
			if _, ok := sqlMap[stringutil.ToSnakeCase(name)]; !ok {
				return fmt.Errorf("cannot found %s in sqlMap from autoTable", name)
			}
			infraRepositoryFile := gen_scaffold.GenerateRepository(name, targetStruct, sqlMap)
			file.JenniferToFile(infraRepositoryFile, cnf.GetInfraRepoFilePath(name))
			return err
		},
	}
	allCmd.PersistentFlags().StringP("pfile", "f", "default", "proto file name")

	return allCmd
}

func NewScaffoldModelCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	modelCmd := &cobra.Command{
		Use:           "model",
		Short:         "generate model struct.",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ctx.IsInsideApp() {
				return errors.New("roche command should be execute inside a rochectl application directory")
			}
			if !cnf.FindToml {
				return errors.New("For using this roche command, please run following command\nrochectl toml\nAnd please edit rochectl.toml According to your project")
			}
			name := util.SnakeToUpperCamel(args[0])
			pfile, err := cmd.PersistentFlags().GetString("pfile")
			if err != nil {
				return err
			}
			if pfile == "default" {
				pfile = name
			}
			pbGoFilePath := cnf.GetPbGoFilePath(pfile)
			targetStruct := ast.FindStruct(name, pbGoFilePath)
			if targetStruct == nil {
				return errors.New("found " + pbGoFilePath + "but not found" + name + " struct")
			}
			infraModelFile := gen_scaffold.GenerateModel(name, targetStruct, ctx.Build.AppName)
			file.JenniferToFile(infraModelFile, cnf.GetInfraModelFilePath(name))
			return err
		},
	}
	modelCmd.PersistentFlags().StringP("pfile", "f", "default", "proto file name")

	return modelCmd
}
