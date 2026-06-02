package casow

import "testing"

func TestNewExpression_shouldStoreCopiedTermsAndConstant_whenGivenTerms(t *testing.T) {
	x := NewVariable()
	y := NewVariable()
	terms := []Term{NewTerm(x, 2.0), NewTerm(y, -3.0)}

	expression := NewExpression(terms, 7.5)
	terms[0] = NewTerm(y, 99.0)

	assertExpression(t, expression, []Term{NewTerm(x, 2.0), NewTerm(y, -3.0)}, 7.5)
}

func TestConstantExpression_shouldStoreConstantWithNoTerms_whenGivenConstant(t *testing.T) {
	expression := ConstantExpression(-4.25)

	assertExpression(t, expression, nil, -4.25)
}

func TestExpressionFromTerm_shouldStoreOneTermAndZeroConstant_whenGivenTerm(t *testing.T) {
	x := NewVariable()
	term := NewTerm(x, 3.5)

	expression := ExpressionFromTerm(term)

	assertExpression(t, expression, []Term{term}, 0.0)
}

func TestExpressionFromTerms_shouldStoreCopiedTermsAndZeroConstant_whenGivenTerms(t *testing.T) {
	x := NewVariable()
	y := NewVariable()
	terms := []Term{NewTerm(x, 1.5), NewTerm(y, 2.5)}

	expression := ExpressionFromTerms(terms)
	terms[1] = NewTerm(x, 99.0)

	assertExpression(t, expression, []Term{NewTerm(x, 1.5), NewTerm(y, 2.5)}, 0.0)
}

func TestExpressionFromVariable_shouldStoreUnitTermAndZeroConstant_whenGivenVariable(t *testing.T) {
	x := NewVariable()

	expression := ExpressionFromVariable(x)

	assertExpression(t, expression, []Term{NewTerm(x, 1.0)}, 0.0)
}

func TestVar_shouldStoreUnitTermAndZeroConstant_whenGivenVariable(t *testing.T) {
	x := NewVariable()

	expression := Var(x)

	assertExpression(t, expression, []Term{NewTerm(x, 1.0)}, 0.0)
}

func TestConst_shouldStoreConstantWithNoTerms_whenGivenConstant(t *testing.T) {
	expression := Const(12.5)

	assertExpression(t, expression, nil, 12.5)
}

func TestTerms_shouldReturnCopy_whenCallerMutatesReturnedSlice(t *testing.T) {
	x := NewVariable()
	y := NewVariable()
	expression := NewExpression([]Term{NewTerm(x, 2.0), NewTerm(y, 3.0)}, 5.0)

	terms := expression.Terms()
	terms[0] = NewTerm(y, 99.0)

	assertExpression(t, expression, []Term{NewTerm(x, 2.0), NewTerm(y, 3.0)}, 5.0)
}

func TestNegate_shouldNegateTermsAndConstant_whenExpressionHasTermsAndConstant(t *testing.T) {
	x := NewVariable()
	y := NewVariable()
	expression := NewExpression([]Term{NewTerm(x, 2.0), NewTerm(y, -3.0)}, 4.0)

	negated := expression.Negate()

	assertExpression(t, negated, []Term{NewTerm(x, -2.0), NewTerm(y, 3.0)}, -4.0)
}

func TestMul_shouldScaleTermsAndConstant_whenGivenMultiplier(t *testing.T) {
	x := NewVariable()
	y := NewVariable()
	expression := NewExpression([]Term{NewTerm(x, 2.0), NewTerm(y, -3.0)}, 4.0)

	scaled := expression.Mul(2.5)

	assertExpression(t, scaled, []Term{NewTerm(x, 5.0), NewTerm(y, -7.5)}, 10.0)
}

func TestDiv_shouldScaleTermsAndConstant_whenGivenDivisor(t *testing.T) {
	x := NewVariable()
	y := NewVariable()
	expression := NewExpression([]Term{NewTerm(x, 6.0), NewTerm(y, -9.0)}, 12.0)

	scaled := expression.Div(3.0)

	assertExpression(t, scaled, []Term{NewTerm(x, 2.0), NewTerm(y, -3.0)}, 4.0)
}

func TestPlusConstant_shouldAdjustOnlyConstant_whenGivenValue(t *testing.T) {
	x := NewVariable()
	expression := NewExpression([]Term{NewTerm(x, 2.0)}, 4.0)

	sum := expression.PlusConstant(1.5)

	assertExpression(t, sum, []Term{NewTerm(x, 2.0)}, 5.5)
}

func TestMinusConstant_shouldAdjustOnlyConstant_whenGivenValue(t *testing.T) {
	x := NewVariable()
	expression := NewExpression([]Term{NewTerm(x, 2.0)}, 4.0)

	difference := expression.MinusConstant(1.5)

	assertExpression(t, difference, []Term{NewTerm(x, 2.0)}, 2.5)
}

func TestPlusExpression_shouldAppendTermsAndAddConstant_whenGivenExpression(t *testing.T) {
	x := NewVariable()
	y := NewVariable()
	left := NewExpression([]Term{NewTerm(x, 2.0)}, 4.0)
	right := NewExpression([]Term{NewTerm(y, 3.0)}, 5.0)

	sum := left.PlusExpression(right)

	assertExpression(t, sum, []Term{NewTerm(x, 2.0), NewTerm(y, 3.0)}, 9.0)
}

func TestMinusExpression_shouldAppendNegatedTermsAndSubtractConstant_whenGivenExpression(t *testing.T) {
	x := NewVariable()
	y := NewVariable()
	left := NewExpression([]Term{NewTerm(x, 2.0)}, 4.0)
	right := NewExpression([]Term{NewTerm(y, 3.0)}, 5.0)

	difference := left.MinusExpression(right)

	assertExpression(t, difference, []Term{NewTerm(x, 2.0), NewTerm(y, -3.0)}, -1.0)
}

func assertExpression(t *testing.T, expression Expression, wantTerms []Term, wantConstant float64) {
	t.Helper()

	if expression.Constant() != wantConstant {
		t.Fatalf("Constant() = %v, want %v", expression.Constant(), wantConstant)
	}

	gotTerms := expression.Terms()
	if len(gotTerms) != len(wantTerms) {
		t.Fatalf("len(Terms()) = %d, want %d", len(gotTerms), len(wantTerms))
	}
	for i, got := range gotTerms {
		want := wantTerms[i]
		if got.Var().ID() != want.Var().ID() || got.Coefficient() != want.Coefficient() {
			t.Fatalf("Terms()[%d] = (%d, %v), want (%d, %v)", i, got.Var().ID(), got.Coefficient(), want.Var().ID(), want.Coefficient())
		}
	}
}
