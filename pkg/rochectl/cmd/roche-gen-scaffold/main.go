package roche_gen_scaffold

import (
	"errors"

	"fmt"
	goAst "go/ast"

	autoTable "github.com/hourglasshoro/auto-table/pkg"
	autoTableFile "github.com/hourglasshoro/auto-table/pkg/file"
	"github.com/izumin5210/grapi/pkg/grapicmd"
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
	scaffoldCmd.AddCommand(NewScaffoldDomainCommand(ctx, cnf))
	scaffoldCmd.AddCommand(NewScaffoldRepositoryCommand(ctx, cnf))
	scaffoldCmd.AddCommand(NewScaffoldMigrationCommand(ctx, cnf))
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
			domainRepositoryFile, usecaseFile := gen_scaffold.GenerateUsecase(name, targetStruct, cnf.ModuleName)

			usecaseProviderName := "New" + name + "Usecase"
			usecaseProviderFile, err := gen_scaffold.GenerateProviderFile(cnf.UsecaseDir, usecaseProviderName, nil)
			if err != nil {
				return err
			}
			file.CreateAndWrite(usecaseProviderFile, cnf.UsecaseDir+"/"+"provider.go")

			file.JenniferToFile(usecaseFile, cnf.GetUsecaseFilePath(name))
			file.JenniferToFile(domainRepositoryFile, cnf.GetDomainRepoFilePath(name))

			repoProviderName := "New" + name + "Repository"
			bMap := map[string]*ast.InterfaceSpec{
				name + "Repository": &ast.InterfaceSpec{
					Name:       "I" + name + "Repository",
					ImportPath: cnf.ModuleName + "/" + cnf.DomainRepoDir,
				},
			}
			// infra/repo層にprovider周りを追加
			infraRepoProviderFile, err := gen_scaffold.GenerateProviderFile(cnf.RepoDir, repoProviderName, bMap)
			if err != nil {
				return err
			}
			file.CreateAndWrite(infraRepoProviderFile, cnf.RepoDir+"/"+"provider.go")

			importPathList := []string{
				cnf.ModuleName + "/" + cnf.RepoDir,
				cnf.ModuleName + "/" + cnf.UsecaseDir,
			}
			wireFile, err := gen_scaffold.GenerateWireFile(cnf.DiDir, importPathList)
			if err != nil {
				return err
			}
			file.CreateAndWrite(wireFile, cnf.DiDir+"/"+"wire.go")

			infraModelFile := gen_scaffold.GenerateModel(name, targetStruct, ctx.Build.AppName)
			file.JenniferToFile(infraModelFile, cnf.GetInfraModelFilePath(name))

			generateEntityAndDomainRepositoryFile(name, targetStruct, cnf)
			generateModelFile(name, targetStruct, ctx.Build.AppName, cnf)

			// TODO: Do refactoring
			generator := autoTable.NewGenerator(ctx.Build.AppName)
			files, err := autoTableFile.GetFiles(&ctx.FS, cnf.InfraModelDir)
			if err != nil {
				return xerrors.Errorf("cannot read infra model files: %w", err)
			}
			sqlMap, err := generator.CreateSQL(files)
			if _, ok := sqlMap[util.CamelToSnake(name)]; !ok {
				return fmt.Errorf("cannot found %s in sqlMap from autoTable", name)
			}
			infraRepositoryFile := gen_scaffold.GenerateRepository(name, targetStruct, sqlMap, cnf.ModuleName)
			file.JenniferToFile(infraRepositoryFile, cnf.GetInfraRepoFilePath(name))
			err = generator.WriteFile(&ctx.FS, cnf.MigrationDir, file.CreateAndWrite)
			return err
		},
	}
	allCmd.PersistentFlags().StringP("pfile", "f", "default", "proto file name")

	return allCmd
}

func NewScaffoldDomainCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	domainCmd := &cobra.Command{
		Use:           "domain",
		Short:         "generate domain method.",
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
			generateEntityAndDomainRepositoryFile(name, targetStruct, cnf)
			return err
		},
	}
	domainCmd.PersistentFlags().StringP("pfile", "f", "default", "proto file name")
	return domainCmd
}

func generateEntityAndDomainRepositoryFile(name string, targetStruct *goAst.StructType, cnf *config.Config) {
	entityFile := gen_scaffold.GenerateEntity(name, targetStruct)
	file.JenniferToFile(entityFile, cnf.GetEntityFilePath(name))
	domainRepositoryFile, usecaseFile := gen_scaffold.GenerateUsecase(name, targetStruct, cnf.ModuleName)
	file.JenniferToFile(usecaseFile, cnf.GetUsecaseFilePath(name))
	file.JenniferToFile(domainRepositoryFile, cnf.GetDomainRepoFilePath(name))
}

func NewScaffoldRepositoryCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	repoCmd := &cobra.Command{
		Use:           "repo",
		Short:         "generate infra repo method.",
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

			// TODO: Do refactoring
			generator := autoTable.NewGenerator(ctx.Build.AppName)
			files, err := autoTableFile.GetFiles(&ctx.FS, cnf.InfraModelDir)
			if err != nil {
				return xerrors.Errorf("cannot read infra model files: %w", err)
			}
			sqlMap, err := generator.CreateSQL(files)
			if _, ok := sqlMap[util.CamelToSnake(name)]; !ok {
				return fmt.Errorf("cannot found %s in sqlMap from autoTable", name)
			}
			infraRepositoryFile := gen_scaffold.GenerateRepository(name, targetStruct, sqlMap, cnf.ModuleName)
			file.JenniferToFile(infraRepositoryFile, cnf.GetInfraRepoFilePath(name))
			return err
		},
	}
	repoCmd.PersistentFlags().StringP("pfile", "f", "default", "proto file name")

	return repoCmd

}

func NewScaffoldMigrationCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:           "migrate",
		Short:         "generate migration file.",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ctx.IsInsideApp() {
				return errors.New("roche command should be execute inside a rochectl application directory")
			}
			if !cnf.FindToml {
				return errors.New("For using this roche command, please run following command\nrochectl toml\nAnd please edit rochectl.toml According to your project")
			}

			return generateMigrationFile(ctx, cnf)
		},
	}
	migrateCmd.PersistentFlags().StringP("pfile", "f", "default", "proto file name")

	return migrateCmd

}

func generateMigrationFile(ctx *grapicmd.Ctx, cnf *config.Config) error {
	generator := autoTable.NewGenerator(ctx.Build.AppName)
	err := generator.WriteFile(&ctx.FS, cnf.MigrationDir, file.CreateAndWrite)
	return err
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
			generateModelFile(name, targetStruct, ctx.Build.AppName, cnf)
			infraModelFile := gen_scaffold.GenerateModel(name, targetStruct, ctx.Build.AppName)
			file.JenniferToFile(infraModelFile, cnf.GetInfraModelFilePath(name))
			return err
		},
	}
	modelCmd.PersistentFlags().StringP("pfile", "f", "default", "proto file name")

	return modelCmd
}

func generateModelFile(name string, targetStruct *goAst.StructType, appName string, cnf *config.Config) {
	infraModelFile := gen_scaffold.GenerateModel(name, targetStruct, appName)
	file.JenniferToFile(infraModelFile, cnf.GetInfraModelFilePath(name))
}
