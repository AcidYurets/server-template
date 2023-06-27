package http

import (
	"context"
	"server-template/internal/modules/domain/entity/dto"
	"server-template/internal/modules/domain/entity/service"
	"server-template/internal/modules/http"
)

type EntityController struct {
	service *service.EntityService
}

func NewEntityController(service *service.EntityService) *EntityController {
	return &EntityController{
		service: service,
	}
}

func InvokeEntityController(controller *EntityController, router http.ApiRouter) {
	router.Post("/entity", controller.Create)
}

func (controller *EntityController) Create(ctx context.Context, entity *dto.EntityCreate) (*dto.Entity, error) {
	createdEntity, err := controller.service.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return createdEntity, nil
}
