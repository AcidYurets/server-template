package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Entity holds the schema definition for the Entity entity.
type Entity struct {
	ent.Schema
}

// Fields of the Entity.
func (Entity) Fields() []ent.Field {
	return []ent.Field{
		field.String("field"),
	}
}

// Edges of the Entity.
func (Entity) Edges() []ent.Edge {
	return []ent.Edge{}
}
