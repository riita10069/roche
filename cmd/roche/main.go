package main

import (
	"fmt"
	roche_custom_method "github.com/riita10069/roche/pkg/rochectl/cmd/roche-custom-method"
	roche_gen_scaffold "github.com/riita10069/roche/pkg/rochectl/cmd/roche-gen-scaffold"
	roche_gen_toml "github.com/riita10069/roche/pkg/rochectl/cmd/roche-gen-toml"
	roche_manifest "github.com/riita10069/roche/pkg/rochectl/cmd/roche-manifest"
	roche_test "github.com/riita10069/roche/pkg/rochectl/cmd/roche-test"
	"github.com/riita10069/roche/pkg/rochectl/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"

	"github.com/izumin5210/clig/pkg/clib"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/cmd"
)

const (
	name     = "roche"
	version  = "v0.5.0"
	tomlFile = "roche.toml"
)

var (
	// set via ldflags
	revision  string
	buildDate string
)


func main() {
	var cnf config.Config

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	grapictx := &grapicmd.Ctx{
		IO:      clib.Stdio(),
		RootDir: cli.RootDir{clib.Path(cwd)},
		Build: clib.Build{
			AppName:   name,
			Version:   version,
			Revision:  revision,
			BuildDate: buildDate,
		},
	}
	command := cmd.NewGrapiCommand(grapictx)
	command.AddCommand(
		roche_test.NewTestCommand(grapictx, &cnf),
		roche_gen_scaffold.NewScaffoldCommand(grapictx, &cnf),
		roche_gen_toml.NewTomlCommand(),
		roche_manifest.NewManifestCommand(grapictx, &cnf),
		roche_custom_method.NewCustomMethodCommand(grapictx, &cnf),
		)

	cobra.OnInitialize(func() {
		if _, err := os.Stat(tomlFile); os.IsNotExist(err) {
			cnf.FindToml = false
			fmt.Println("cannot find roche.toml")
		} else {
			cnf.FindToml = true
			viper.SetConfigFile(tomlFile)

			// 設定ファイルを読み込む
			if err := viper.ReadInConfig(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// 設定ファイルの内容を構造体にコピーする
			if err := viper.Unmarshal(&cnf); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	})

	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
