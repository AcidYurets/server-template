package db

import (
	"go.uber.org/fx"
)

var Module = fx.Module("db",
	fx.Provide(NewDBClient),
	fx.Invoke(InvokeDBClient))
