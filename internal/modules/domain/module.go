package domain

import (
	"go.uber.org/fx"
	"server-template/internal/modules/domain/entity"
)

var (
	Module = fx.Module("domain",
		entity.Module,
	)
)
