package casow

import "testing"

func TestNewTermStoresVariableAndCoefficient(t *testing.T) {
	variable := NewVariable()

	term := NewTerm(variable, 2.5)

	if term.Var() != variable {
		t.Fatalf("NewTerm variable = %v, want %v", term.Var(), variable)
	}
	if term.Coefficient() != 2.5 {
		t.Fatalf("NewTerm coefficient = %v, want 2.5", term.Coefficient())
	}
}

func TestTermFromVariableUsesUnitCoefficient(t *testing.T) {
	variable := NewVariable()

	term := TermFromVariable(variable)

	if term.Var() != variable {
		t.Fatalf("TermFromVariable variable = %v, want %v", term.Var(), variable)
	}
	if term.Coefficient() != 1.0 {
		t.Fatalf("TermFromVariable coefficient = %v, want 1.0", term.Coefficient())
	}
}

func TestTermAccessorsReturnStoredValues(t *testing.T) {
	variable := NewVariable()
	term := NewTerm(variable, -3.25)

	if term.Var() != variable {
		t.Fatalf("Var() = %v, want %v", term.Var(), variable)
	}
	if term.Coefficient() != -3.25 {
		t.Fatalf("Coefficient() = %v, want -3.25", term.Coefficient())
	}
}

func TestTermArithmeticReturnsTransformedTerms(t *testing.T) {
	variable := NewVariable()

	tests := []struct {
		name     string
		got      Term
		expected float64
	}{
		{name: "mul", got: NewTerm(variable, 2.0).Mul(3.5), expected: 7.0},
		{name: "div", got: NewTerm(variable, 9.0).Div(3.0), expected: 3.0},
		{name: "negate", got: NewTerm(variable, 4.5).Negate(), expected: -4.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got.Var() != variable {
				t.Fatalf("arithmetic variable = %v, want %v", tt.got.Var(), variable)
			}
			if tt.got.Coefficient() != tt.expected {
				t.Fatalf("arithmetic coefficient = %v, want %v", tt.got.Coefficient(), tt.expected)
			}
		})
	}
}

func TestTermChainedArithmeticPreservesVariableIdentity(t *testing.T) {
	variable := NewVariable()

	term := NewTerm(variable, 6.0).Mul(2.0).Div(3.0).Negate()

	if term.Var() != variable {
		t.Fatalf("chained arithmetic variable = %v, want %v", term.Var(), variable)
	}
	if term.Coefficient() != -4.0 {
		t.Fatalf("chained arithmetic coefficient = %v, want -4.0", term.Coefficient())
	}
}
