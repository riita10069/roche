package main

import (
	"fmt"
	roche_gen_scaffold "github.com/riita10069/roche/pkg/roche/cmd/roche-gen-scaffold"
	roche_gen_toml "github.com/riita10069/roche/pkg/roche/cmd/roche-gen-toml"
	roche_test "github.com/riita10069/roche/pkg/roche/cmd/roche-test"
	"github.com/riita10069/roche/pkg/roche/config"
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
		roche_gen_toml.NewTomlCommand(grapictx),
		)

	cobra.OnInitialize(func() {
		_, err := os.Stat(tomlFile)
		if err != nil {
			roche_gen_toml.MakeRocheToml()
		}

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
	})

	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
