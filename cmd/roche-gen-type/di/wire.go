//+build wireinject

package di

import (
	"github.com/google/wire"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/gencmd"
	"github.com/izumin5210/grapi/pkg/protoc"
)

func NewApp(*gencmd.Command) (*App, error) {
	wire.Build(
		App{},
		gencmd.Set,
		cli.UIInstance,
		protoc.WrapperSet,
	)
	return nil, nil
}
