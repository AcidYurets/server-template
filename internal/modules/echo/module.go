package echo

import (
	"go.uber.org/fx"
)

var Module = fx.Module("echo",
	fx.Provide(NewEcho),
	fx.Invoke(InvokeEcho))
