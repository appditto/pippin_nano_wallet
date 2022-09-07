package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// AdhocAccount holds the schema definition for the AdhocAccount entity.
type AdhocAccount struct {
	ent.Schema
}

func (AdhocAccount) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "adhoc_accounts"},
	}
}

// Fields of the AdhocAccount.
func (AdhocAccount) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("wallet_id", uuid.UUID{}),
		field.String("address").MaxLen(65),
		field.String("private_key").MaxLen(128),
		field.Bool("work").Default(true),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the AdhocAccount.
func (AdhocAccount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("wallet", Wallet.Type).
			Ref("adhoc_accounts").
			Field("wallet_id").
			Required().
			Unique(),
		edge.To("blocks", Block.Type),
	}
}

func (AdhocAccount) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("wallet_id"),
		index.Fields("address", "account_index").Unique(),
	}
}
