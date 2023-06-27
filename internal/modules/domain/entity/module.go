package entity

import (
	"go.uber.org/fx"
	"server-template/internal/modules/domain/entity/http"
	"server-template/internal/modules/domain/entity/repo"
	"server-template/internal/modules/domain/entity/service"
)

var (
	Module = fx.Module("entity",
		http.Module,
		service.Module,
		repo.Module,

		fx.Provide(
			fx.Annotate(
				func(r *repo.EntityRepo) *repo.EntityRepo { return r },
				fx.As(new(service.IEntityRepo)),
			),
		),
	)
)
