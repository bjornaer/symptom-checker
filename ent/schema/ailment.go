package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type SymptomDetails struct {
	Name      string `json:"name"`
	HPOId     string `json:"HPOId"`
	Frequency string `json:"frequency"`
}

// Ailment holds the schema definition for the Ailment entity.
type Ailment struct {
	ent.Schema
}

// Fields of the Ailment.
func (Ailment) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive(),
		field.String("name").
			Default("unknown"),
		field.JSON("symptoms", map[string]SymptomDetails{}).
			Default(map[string]SymptomDetails{}),
		field.Strings("hpos").
			Default([]string{}),
		field.String("expert").
			Default(""),
	}
}

// Edges of the Ailment.
func (Ailment) Edges() []ent.Edge {
	return nil
}
