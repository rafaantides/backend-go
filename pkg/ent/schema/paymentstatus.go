package schema

import (
	"backend-go/pkg/mixins"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type PaymentStatus struct {
	ent.Schema
}

func (PaymentStatus) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.UUIDMixin{},
		mixins.TimestampsMixin{},
	}
}

func (PaymentStatus) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(100),
		field.String("description").Optional().Nillable(),
	}
}
