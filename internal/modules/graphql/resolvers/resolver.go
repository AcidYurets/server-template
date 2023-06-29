package resolvers

import entity_srv "server-template/internal/modules/domain/entity/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	entityService *entity_srv.EntityService
}

func NewResolver(entityService *entity_srv.EntityService) *Resolver {
	return &Resolver{entityService: entityService}
}
