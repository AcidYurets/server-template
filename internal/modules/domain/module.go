package domain

import (
	"go.uber.org/fx"
	"server-template/internal/modules/domain/entity"
)

var (
	Module = fx.Options(
		entity.Module,
	)
	Invokables = fx.Options(
		entity.Invokables,
	)
)
