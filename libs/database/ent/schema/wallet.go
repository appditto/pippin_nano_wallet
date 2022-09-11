package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Wallet holds the schema definition for the Wallet entity.
type Wallet struct {
	ent.Schema
}

// Annotations of the Wallet.
func (Wallet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "wallets"},
	}
}

// Fields of the Wallet.
func (Wallet) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		// Large enough to store encrypted keys, which have more bits
		field.String("seed").MaxLen(512).Unique(),
		field.String("representative").MaxLen(65).Nillable().Optional(),
		field.Bool("encrypted").Default(false),
		field.Bool("work").Default(true),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the Wallet.
func (Wallet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("accounts", Account.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("adhoc_accounts", AdhocAccount.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
