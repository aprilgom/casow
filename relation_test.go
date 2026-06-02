package casow

import "testing"

func TestRelationalOperatorString_shouldReturnSymbol_whenKnownOperator(t *testing.T) {
	tests := []struct {
		name     string
		operator RelationalOperator
		expected string
	}{
		{name: "less or equal", operator: LessOrEqual, expected: "<="},
		{name: "equal", operator: Equal, expected: "=="},
		{name: "greater or equal", operator: GreaterOrEqual, expected: ">="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operator.String(); got != tt.expected {
				t.Fatalf("%v.String() = %q, want %q", tt.operator, got, tt.expected)
			}
		})
	}
}

func TestWeightedRelationConstructors_shouldStoreOperatorAndStrength_whenGivenStrength(t *testing.T) {
	tests := []struct {
		name     string
		relation WeightedRelation
		operator RelationalOperator
		strength Strength
	}{
		{name: "EQ", relation: EQ(Required), operator: Equal, strength: Required},
		{name: "LE", relation: LE(Strong), operator: LessOrEqual, strength: Strong},
		{name: "GE", relation: GE(Weak), operator: GreaterOrEqual, strength: Weak},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.relation.Operator(); got != tt.operator {
				t.Fatalf("Operator() = %v, want %v", got, tt.operator)
			}
			if got := tt.relation.Strength(); got != tt.strength {
				t.Fatalf("Strength() = %v, want %v", got, tt.strength)
			}
		})
	}
}

func TestWeightedRelationConstructors_shouldExposeDataForNewConstraint_whenUsedAsPublicHelpers(t *testing.T) {
	x := NewVariable()
	tests := []struct {
		name     string
		relation WeightedRelation
		operator RelationalOperator
		strength Strength
	}{
		{name: "EQ", relation: EQ(Required), operator: Equal, strength: Required},
		{name: "LE", relation: LE(Medium), operator: LessOrEqual, strength: Medium},
		{name: "GE", relation: GE(Strong), operator: GreaterOrEqual, strength: Strong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := NewConstraint(x, tt.relation.Operator(), 10.0, tt.relation.Strength())

			if got := constraint.Operator(); got != tt.operator {
				t.Fatalf("Operator() = %v, want %v", got, tt.operator)
			}
			if got := constraint.Strength(); got != tt.strength {
				t.Fatalf("Strength() = %v, want %v", got, tt.strength)
			}
			assertExpression(t, constraint.Expression(), []Term{NewTerm(x, 1.0)}, -10.0)
		})
	}
}
