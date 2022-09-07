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

// Block holds the schema definition for the Block entity.
type Block struct {
	ent.Schema
}

// Set annotations for the Block.
func (Block) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "blocks"},
	}
}

// Fields of the Block.
func (Block) Fields() []ent.Field {
	return []ent.Field{
		// account id from accounts
		field.UUID("account_id", uuid.UUID{}),
		field.UUID("adhoc_account_id", uuid.UUID{}),
		field.String("block_hash").MaxLen(64).Unique(),
		// TODO use a proper struct, not map[string]interface{}
		field.JSON("block", map[string]interface{}{}),
		field.String("send_id").MaxLen(128).Nillable().Optional(),
		field.String("subtype").MaxLen(10),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Block.
func (Block) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).
			Ref("blocks").
			Field("account_id").
			Required().
			Unique(),
		edge.From("adhoc_account", AdhocAccount.Type).
			Ref("blocks").
			Field("adhoc_account_id").
			Required().
			Unique(),
	}
}

// Indexes of the Block.
func (Block) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("account_id"),
		index.Fields("adhoc_account_id"),
		index.Fields("send_id"),
		index.Fields("account_id", "send_id").Unique(),
	}
}
