package http

import (
	"go.uber.org/fx"
)

var (
	Module     = fx.Provide(NewEntityController)
	Invokables = fx.Invoke(InvokeEntityController)
)
