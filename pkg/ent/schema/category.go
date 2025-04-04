package schema

import (
	"backend-go/pkg/mixins"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Category struct {
	ent.Schema
}

func (Category) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.UUIDMixin{},
		mixins.TimestampsMixin{},
	}
}

func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(255),
		field.String("description").Optional().Nillable(),
	}
}
