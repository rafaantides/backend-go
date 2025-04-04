package schema

import (
	"backend-go/pkg/mixins"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Debt struct {
	ent.Schema
}

func (Debt) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.UUIDMixin{},
		mixins.TimestampsMixin{},
		mixins.MoneyMixin{Name: "amount"},
	}
}

func (Debt) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").MaxLen(255),
		field.Time("purchase_date"),
		field.Time("due_date"),
	}
}

func (Debt) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("invoice", Invoice.Type).Unique().StorageKey(edge.Column("invoice_id")),
		edge.To("category", Category.Type).Unique().StorageKey(edge.Column("category_id")),
		edge.To("status", PaymentStatus.Type).Unique().StorageKey(edge.Column("status_id")),
	}
}
