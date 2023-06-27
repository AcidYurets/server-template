package repo

import (
	"context"
	"server-template/internal/modules/db/ent"
	"server-template/internal/modules/domain/entity/dto"
)

type EntityRepo struct {
	client *ent.Client
}

func NewEntityRepo(client *ent.Client) *EntityRepo {
	return &EntityRepo{
		client: client,
	}
}

// Create создание записи в БД
func (r *EntityRepo) Create(ctx context.Context, entity *dto.EntityCreate) (*dto.Entity, error) {
	newEntity, err := r.client.Entity.Create().
		SetField(entity.Field).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return ToEntityDTO(newEntity), nil
}

func ToEntityDTO(model *ent.Entity) *dto.Entity {
	if model == nil {
		return nil
	}

	return &dto.Entity{
		Id:    model.ID,
		Field: model.Field,
	}
}

func ToEntityDTOs(models ent.Entities) dto.Entities {
	if models == nil {
		return nil
	}
	dtms := make(dto.Entities, len(models))
	for i := range models {
		dtms[i] = ToEntityDTO(models[i])
	}
	return dtms
}
