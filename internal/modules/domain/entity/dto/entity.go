package dto

// Entity модель сущности
type Entity struct {
	Id    int    // ID сущности
	Filed string // Поле сущности
}

type Entities []*Entity
