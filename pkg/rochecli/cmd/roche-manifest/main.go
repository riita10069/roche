package roche_manifest

import (
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/riita10069/roche/pkg/roche/config"
	gen_manifest "github.com/riita10069/roche/pkg/roche/gen-manifest"
	"github.com/spf13/cobra"
)

func NewManifestCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	manifestCmd := &cobra.Command{
		Use:           "manifest",
		Short:         "k8s manifest generator",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	manifestCmd.AddCommand(NewAppCommand(ctx, cnf))
	return manifestCmd
}

func NewAppCommand(ctx *grapicmd.Ctx, cnf *config.Config) *cobra.Command {
	manifestCmd := &cobra.Command{
		Use:           "app",
		Short:         "k8s manifest generator to make an app",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			err := gen_manifest.GenerateDeployment(name, cnf)
			if err != nil {
				return err
			}
			err = gen_manifest.GenerateHpa(name, cnf)
			if err != nil {
				return err
			}
			err = gen_manifest.GenerateNodePort(name, cnf)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return manifestCmd
}
