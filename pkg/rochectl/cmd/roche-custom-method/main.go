package roche_custom_method

import (
	"errors"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/riita10069/roche/pkg/rochectl/ast"
	"github.com/riita10069/roche/pkg/rochectl/config"
	custom_method "github.com/riita10069/roche/pkg/rochectl/custom-method"
	"github.com/riita10069/roche/pkg/rochectl/file"
	"github.com/riita10069/roche/pkg/util"
	"github.com/spf13/cobra"
)


func NewCustomMethodCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	testCmd := &cobra.Command{
		Use:           "custom NAME",
		Short:         "Add unimplemented proto methods to ServerDir.",
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
			wantToFindInterfaceName := name + "ServiceServer"
			targetInterface := ast.FindInterface(wantToFindInterfaceName, pbGoFilePath)
			if targetInterface == nil {
				return errors.New("cannot find interface at " + pbGoFilePath)
			}

			handlerFilePath := cnf.ServerDir + "/" + util.CamelToSnake(name) + ".go"
			handlerFuncList := ast.FindFunc(handlerFilePath)
			if handlerFuncList == nil {
				return errors.New("cannot find function at " + handlerFilePath)
			}

			customMethodList := custom_method.GetList(targetInterface, handlerFuncList, name)
			if len(customMethodList) == 0 {
				return errors.New("cannot make custom method list")
			}

			for _, customMethod := range customMethodList {
				file.Append("\n" + customMethod, handlerFilePath)
			}
			return nil
		},
	}
	return testCmd
}
