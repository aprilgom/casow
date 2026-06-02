package casow

import "testing"

func TestUsageTwoHorizontalBoxes_shouldUpdateLayout_whenWindowWidthChanges(t *testing.T) {
	type element struct {
		left  Variable
		right Variable
	}

	variable := func(v Variable) Expression {
		return ExpressionFromVariable(v)
	}
	addConstraint := func(solver *Solver, lhs Expression, op RelationalOperator, rhs Expression, strength Strength) {
		t.Helper()
		if err := solver.AddConstraint(NewConstraint(lhs, op, rhs, strength)); err != nil {
			t.Fatalf("AddConstraint error = %v, want nil", err)
		}
	}
	assertEqual := func(name string, got, want float64) {
		t.Helper()
		if got != want {
			t.Fatalf("%s = %v, want %v", name, got, want)
		}
	}

	valueOf, updateValues := newValues()
	windowWidth := NewVariable()
	box1 := element{left: NewVariable(), right: NewVariable()}
	box2 := element{left: NewVariable(), right: NewVariable()}
	solver := NewSolver()

	addConstraint(solver, variable(windowWidth), GreaterOrEqual, ConstantExpression(0), Required)
	addConstraint(solver, variable(box1.left), Equal, ConstantExpression(0), Required)
	addConstraint(solver, variable(box2.right), Equal, variable(windowWidth), Required)
	addConstraint(solver, variable(box2.left), GreaterOrEqual, variable(box1.right), Required)
	addConstraint(solver, variable(box1.left), LessOrEqual, variable(box1.right), Required)
	addConstraint(solver, variable(box2.left), LessOrEqual, variable(box2.right), Required)
	addConstraint(solver, variable(box1.right).MinusExpression(variable(box1.left)), Equal, ConstantExpression(50), Weak)
	addConstraint(solver, variable(box2.right).MinusExpression(variable(box2.left)), Equal, ConstantExpression(100), Weak)

	if err := solver.AddEditVariable(windowWidth, Strong); err != nil {
		t.Fatalf("AddEditVariable(windowWidth) error = %v, want nil", err)
	}
	if err := solver.SuggestValue(windowWidth, 300); err != nil {
		t.Fatalf("SuggestValue(windowWidth, 300) error = %v, want nil", err)
	}
	updateValues(solver.FetchChanges())

	assertEqual("windowWidth after 300 suggestion", valueOf(windowWidth), 300)
	assertEqual("box1.right after 300 suggestion", valueOf(box1.right), 50)
	assertEqual("box2.left after 300 suggestion", valueOf(box2.left), 200)
	assertEqual("box2.right after 300 suggestion", valueOf(box2.right), 300)

	if err := solver.SuggestValue(windowWidth, 75); err != nil {
		t.Fatalf("SuggestValue(windowWidth, 75) error = %v, want nil", err)
	}
	updateValues(solver.FetchChanges())

	assertEqual("windowWidth after 75 suggestion", valueOf(windowWidth), 75)
	assertEqual("box2.right after 75 suggestion", valueOf(box2.right), 75)
	if valueOf(box1.left) != 0 {
		t.Fatalf("box1.left = %v, want 0", valueOf(box1.left))
	}
	if valueOf(box2.left) < valueOf(box1.right) {
		t.Fatalf("boxes overlap: box2.left = %v, box1.right = %v", valueOf(box2.left), valueOf(box1.right))
	}
	if valueOf(box1.left) > valueOf(box1.right) {
		t.Fatalf("box1 has negative width: left = %v, right = %v", valueOf(box1.left), valueOf(box1.right))
	}
	if valueOf(box2.left) > valueOf(box2.right) {
		t.Fatalf("box2 has negative width: left = %v, right = %v", valueOf(box2.left), valueOf(box2.right))
	}
	if valueOf(windowWidth) < 0 {
		t.Fatalf("windowWidth = %v, want nonnegative", valueOf(windowWidth))
	}
}
