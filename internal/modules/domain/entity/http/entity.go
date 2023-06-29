package http

import (
	"context"
	"server-template/internal/modules/domain/entity/dto"
	"server-template/internal/modules/domain/entity/service"
	"server-template/internal/modules/http"
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

func InvokeEntityController(controller *EntityController, router http.ApiRouter) {
	router.Get("/entity/:id", controller.GetById)
	router.Get("/entity", controller.List)
	router.Post("/entity", controller.Create)
	router.Put("/entity/:id", controller.Update)
	router.Delete("/entity/:id", controller.Delete)
}

func (controller *EntityController) GetById(ctx context.Context, params routers.Params) (*dto.Entity, error) {
	id, err := params.GetId()
	if err != nil {
		return nil, err
	}

	entity, err := controller.service.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (controller *EntityController) List(ctx context.Context) ([]*dto.Entity, error) {
	entities, err := controller.service.List(ctx)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (controller *EntityController) Create(ctx context.Context, entity *dto.EntityCreate) (*dto.Entity, error) {
	createdEntity, err := controller.service.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return createdEntity, nil
}

func (controller *EntityController) Update(ctx context.Context, params routers.Params, entity *dto.EntityUpdate) (*dto.Entity, error) {
	id, err := params.GetId()
	if err != nil {
		return nil, err
	}

	updatedEntity, err := controller.service.Update(ctx, id, entity)
	if err != nil {
		return nil, err
	}

	return updatedEntity, nil
}

func (controller *EntityController) Delete(ctx context.Context, params routers.Params) error {
	id, err := params.GetId()
	if err != nil {
		return err
	}

	err = controller.service.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
