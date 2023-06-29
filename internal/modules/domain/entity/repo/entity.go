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

// GetById получение записи
func (r *EntityRepo) GetById(ctx context.Context, id int) (*dto.Entity, error) {
	entity, err := r.client.Entity.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return ToEntityDTO(entity), nil
}

// List получение списка записей
func (r *EntityRepo) List(ctx context.Context) ([]*dto.Entity, error) {
	entities, err := r.client.Entity.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	return ToEntityDTOs(entities), nil
}

// Create создание записи
func (r *EntityRepo) Create(ctx context.Context, entity *dto.EntityCreate) (*dto.Entity, error) {
	newEntity, err := r.client.Entity.Create().
		SetField(entity.Field).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return ToEntityDTO(newEntity), nil
}

// Update обновление записи
func (r *EntityRepo) Update(ctx context.Context, id int, entity *dto.EntityUpdate) (*dto.Entity, error) {
	updEntity, err := r.client.Entity.UpdateOneID(id).
		SetField(entity.Field).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return ToEntityDTO(updEntity), nil
}

// Delete удаление записи
func (r *EntityRepo) Delete(ctx context.Context, id int) error {
	err := r.client.Entity.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
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
