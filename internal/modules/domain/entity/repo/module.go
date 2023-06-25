package repo

import (
	"go.uber.org/fx"
)

var (
	Module     = fx.Provide(NewEntityRepo)
	Invokables = fx.Invoke()
)
