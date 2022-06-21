package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Symptom holds the schema definition for the Symptom entity.
type Symptom struct {
	ent.Schema
}

// Fields of the Symptom.
func (Symptom) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Default("unknown"),
		field.String("hpo").
			Unique().
			NotEmpty(),
	}
}

// Edges of the Symptom.
func (Symptom) Edges() []ent.Edge {
	return nil
}
