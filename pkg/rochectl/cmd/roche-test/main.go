package roche_test

import (
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/riita10069/roche/pkg/rochectl/config"
	"github.com/spf13/cobra"
)


func NewTestCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	testCmd := &cobra.Command{
		Use:           "test NAME",
		Short:         "to give it a shot cobra, viper and so on in order to development.",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return testCmd
}
