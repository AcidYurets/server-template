package http

import (
	"context"
	"server-template/internal/modules/domain/entity/dto"
	"server-template/internal/modules/domain/entity/service"
	"server-template/internal/pkg/routers"
)

type EntityController struct {
	service *service.EntityService
}

func NewEntityController(service *service.EntityService) *EntityController {
	return &EntityController{
		service: service,
	}
}

func InvokeEntityController(controller *EntityController, router routers.Router) {
	router.Post("/entity", controller.Create)
}

func (controller *EntityController) Create(ctx context.Context, entity *dto.Entity) (*dto.Entity, error) {
	entity, err := controller.service.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return entity, nil
}
