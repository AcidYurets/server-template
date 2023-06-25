package echo

import (
	"go.uber.org/fx"
)

var (
	Module     = fx.Provide(NewEcho)
	Invokables = fx.Invoke(InvokeEcho)
)
