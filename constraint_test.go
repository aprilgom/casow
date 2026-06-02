package casow

import "testing"

func TestNewConstraint_shouldNotEqualZeroValueConstraint(t *testing.T) {
	var zero Constraint
	nextConstraintID.Store(0)

	created := NewConstraint(Var(NewVariable()), Equal, Const(1), Required)

	if created == zero {
		t.Fatalf("NewConstraint() = zero value constraint: %v", created)
	}
	if created.id == 0 {
		t.Fatalf("NewConstraint() ID = 0, want non-zero")
	}
}

func TestNewConstraint_shouldStoreCanonicalExpressionOperatorAndStrength_whenGivenExpressions(t *testing.T) {
	x := NewVariable()
	y := NewVariable()
	lhs := NewExpression([]Term{NewTerm(x, 2.0)}, 10.0)
	rhs := NewExpression([]Term{NewTerm(y, 3.0)}, 4.0)

	constraint := NewConstraint(lhs, GreaterOrEqual, rhs, Strong)

	assertExpression(t, constraint.Expression(), []Term{NewTerm(x, 2.0), NewTerm(y, -3.0)}, 6.0)
	if got := constraint.Operator(); got != GreaterOrEqual {
		t.Fatalf("Operator() = %v, want %v", got, GreaterOrEqual)
	}
	if got := constraint.Strength(); got != Strong {
		t.Fatalf("Strength() = %v, want %v", got, Strong)
	}
}

func TestConstraintIdentity_shouldCompareByHandle_whenCopiedOrSeparatelyCreated(t *testing.T) {
	x := NewVariable()
	lhs := ExpressionFromVariable(x)
	rhs := ConstantExpression(10.0)

	original := NewConstraint(lhs, Equal, rhs, Required)
	copied := original
	separate := NewConstraint(lhs, Equal, rhs, Required)

	if copied != original {
		t.Fatalf("copied constraint should compare equal: got %v, want %v", copied, original)
	}
	if separate == original {
		t.Fatalf("separately created identical constraints should not compare equal: got %v and %v", separate, original)
	}
}

func TestConstraintCanBeUsedAsMapKey(t *testing.T) {
	x := NewVariable()
	constraint := NewConstraint(ExpressionFromVariable(x), LessOrEqual, ConstantExpression(100.0), Medium)
	values := map[Constraint]string{
		constraint: "width cap",
	}

	if got := values[constraint]; got != "width cap" {
		t.Fatalf("map lookup by constraint key = %q, want %q", got, "width cap")
	}
}
