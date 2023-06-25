package http

import "go.uber.org/fx"

var (
	Module     = fx.Provide(NewHTTP)
	Invokables = fx.Invoke()
)
