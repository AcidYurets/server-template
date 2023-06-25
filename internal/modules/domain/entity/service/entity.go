package service

import (
	"context"
	"server-template/internal/modules/domain/entity/dto"
)

type IEntityRepo interface {
	Create(ctx context.Context, entity *dto.Entity) (*dto.Entity, error)
}

type EntityService struct {
	repo IEntityRepo
}

func NewEntityService(repo IEntityRepo) *EntityService {
	return &EntityService{
		repo: repo,
	}
}

// Create создание записи в БД
func (r *EntityService) Create(ctx context.Context, entity *dto.Entity) (*dto.Entity, error) {
	return r.repo.Create(ctx, entity)
}
