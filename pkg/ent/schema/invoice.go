package schema

import (
	"backend-go/pkg/mixins"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Invoice struct {
	ent.Schema
}

func (Invoice) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.UUIDMixin{},
		mixins.TimestampsMixin{},
		mixins.MoneyMixin{Name: "amount"},
	}
}

func (Invoice) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").MaxLen(255),
		field.Time("issue_date"),
		field.Time("due_date"),
	}
}

func (Invoice) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("status", PaymentStatus.Type).Unique().StorageKey(edge.Column("status_id")),
	}
}
