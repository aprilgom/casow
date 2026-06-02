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

func TestTermArithmeticCompositions_shouldMatchUpstreamOperatorTables(t *testing.T) {
	left := NewVariable()
	right := NewVariable()
	leftTerm := TermFromVariable(left)
	rightTerm := TermFromVariable(right)

	tests := []struct {
		name         string
		got          Expression
		wantTerms    []Term
		wantConstant float64
	}{
		{name: "add constant", got: ExpressionFromTerm(leftTerm).PlusConstant(2.0), wantTerms: []Term{leftTerm}, wantConstant: 2.0},
		{name: "constant add term", got: Const(2.0).PlusExpression(ExpressionFromTerm(leftTerm)), wantTerms: []Term{leftTerm}, wantConstant: 2.0},
		{name: "add term", got: ExpressionFromTerm(leftTerm).PlusExpression(ExpressionFromTerm(rightTerm)), wantTerms: []Term{leftTerm, rightTerm}, wantConstant: 0.0},
		{name: "add expression", got: ExpressionFromTerm(leftTerm).PlusExpression(NewExpression([]Term{rightTerm}, 1.0)), wantTerms: []Term{leftTerm, rightTerm}, wantConstant: 1.0},
		{name: "sub constant", got: ExpressionFromTerm(leftTerm).MinusConstant(2.0), wantTerms: []Term{leftTerm}, wantConstant: -2.0},
		{name: "constant sub term", got: Const(2.0).MinusExpression(ExpressionFromTerm(leftTerm)), wantTerms: []Term{leftTerm.Negate()}, wantConstant: 2.0},
		{name: "sub term", got: ExpressionFromTerm(leftTerm).MinusExpression(ExpressionFromTerm(rightTerm)), wantTerms: []Term{leftTerm, rightTerm.Negate()}, wantConstant: 0.0},
		{name: "sub expression", got: ExpressionFromTerm(leftTerm).MinusExpression(NewExpression([]Term{rightTerm}, 1.0)), wantTerms: []Term{leftTerm, rightTerm.Negate()}, wantConstant: -1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertExpression(t, tt.got, tt.wantTerms, tt.wantConstant)
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
