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

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Annotations of the Account.
func (Account) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "accounts"},
	}
}

// Fields of the Account.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.UUID("wallet_id", uuid.UUID{}),
		field.String("address").MaxLen(65),
		field.Int("account_index"),
		field.Bool("work").Default(true),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the Account.
func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("wallet", Wallet.Type).
			Ref("accounts").
			Field("wallet_id").
			Required().
			Unique(),
		edge.To("blocks", Block.Type),
	}
}

// Indexes of the Wallet.
func (Account) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("wallet_id"),
		index.Fields("wallet_id", "address", "account_index").Unique(),
	}
}
