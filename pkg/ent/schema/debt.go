package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Debt struct {
	ent.Schema
}

func (Debt) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("invoice_id", uuid.UUID{}).Optional().Nillable(),
		field.String("title"),
		field.UUID("category_id", uuid.UUID{}).Optional().Nillable(),
		field.Float("amount"),
		field.Time("purchase_date"),
		field.Time("due_date"),
		field.UUID("status_id", uuid.UUID{}).Optional().Nillable(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Debt) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("invoice", Invoice.Type).Unique().Field("invoice_id"),
		edge.To("category", Category.Type).Unique().Field("category_id"),
		edge.To("status", PaymentStatus.Type).Unique().Field("status_id"),
	}
}
