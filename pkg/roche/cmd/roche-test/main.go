package roche_test

import (
	"fmt"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/riita10069/roche/pkg/roche/config"
	"github.com/spf13/cobra"
)


func NewTestCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	testCmd := &cobra.Command{
		Use:           "test",
		Short:         "test cobra, viper and so on.",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("-----toml ctx-------")
			fmt.Printf("(%%#v) %#v\n", *cnf)
			return nil
		},
	}
	return testCmd
}
