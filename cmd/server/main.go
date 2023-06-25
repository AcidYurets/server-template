package main

import (
	"go.uber.org/fx"
	"server-template/internal/modules"
)

func main() {
	fx.New(
		modules.AppModule,
		modules.AppInvokables,
	).Run()
}
