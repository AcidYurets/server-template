package service

import (
	"go.uber.org/fx"
)

var (
	Module     = fx.Provide(NewEntityService)
	Invokables = fx.Invoke()
)
