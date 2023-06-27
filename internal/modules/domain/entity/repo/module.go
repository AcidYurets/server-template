package repo

import (
	"go.uber.org/fx"
)

var Module = fx.Module("repo",
	fx.Provide(NewEntityRepo))
