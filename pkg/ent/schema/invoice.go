package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Invoice struct {
	ent.Schema
}

func (Invoice) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.String("title"),
		field.Float("amount"),
		field.Time("issue_date"),
		field.Time("due_date"),
		field.UUID("status_id", uuid.UUID{}).Optional().Nillable(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Invoice.
func (Invoice) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("status", PaymentStatus.Type).Unique().Field("status_id"),
	}
}
