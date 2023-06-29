package service

import (
	"context"
	"server-template/internal/modules/domain/entity/dto"
)

type IEntityRepo interface {
	GetById(ctx context.Context, id int) (*dto.Entity, error)
	List(ctx context.Context) ([]*dto.Entity, error)
	Create(ctx context.Context, entity *dto.EntityCreate) (*dto.Entity, error)
	Update(ctx context.Context, id int, entity *dto.EntityUpdate) (*dto.Entity, error)
	Delete(ctx context.Context, id int) error
}

type EntityService struct {
	repo IEntityRepo
}

func NewEntityService(repo IEntityRepo) *EntityService {
	return &EntityService{
		repo: repo,
	}
}

// GetById получение записи
func (r *EntityService) GetById(ctx context.Context, id int) (*dto.Entity, error) {
	entity, err := r.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// List получение списка записей
func (r *EntityService) List(ctx context.Context) ([]*dto.Entity, error) {
	entities, err := r.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

// Create создание записи в БД
func (r *EntityService) Create(ctx context.Context, entity *dto.EntityCreate) (*dto.Entity, error) {
	newEntity, err := r.repo.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return newEntity, nil
}

// Update обновление записи
func (r *EntityService) Update(ctx context.Context, id int, entity *dto.EntityUpdate) (*dto.Entity, error) {
	newEntity, err := r.repo.Update(ctx, id, entity)
	if err != nil {
		return nil, err
	}

	return newEntity, nil
}

// Delete удаление записи
func (r *EntityService) Delete(ctx context.Context, id int) error {
	err := r.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
