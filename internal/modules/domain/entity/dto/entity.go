package dto

// Entity модель сущности
type Entity struct {
	Id    int    // ID сущности
	Field string // Поле сущности
}

type Entities []*Entity

type EntityCreate struct {
	Field string // Поле сущности
}

type EntityUpdate struct {
	Field string // Поле сущности
}
