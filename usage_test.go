package casow

import (
	"math"
	"testing"
)

func TestUsageTwoHorizontalBoxes_shouldUpdateLayout_whenWindowWidthChanges(t *testing.T) {
	type element struct {
		left  Variable
		right Variable
	}

	assertEqual := func(name string, got, want float64) {
		t.Helper()
		if got != want {
			t.Fatalf("%s = %v, want %v", name, got, want)
		}
	}
	assertNear := func(name string, got, want float64) {
		t.Helper()
		if math.Abs(got-want) > 1e-9 {
			t.Fatalf("%s = %v, want %v", name, got, want)
		}
	}

	valueOf, updateValues := newValues()
	windowWidth := NewVariable()
	box1 := element{left: NewVariable(), right: NewVariable()}
	box2 := element{left: NewVariable(), right: NewVariable()}
	solver := NewSolver()

	if err := solver.AddConstraints(
		NewConstraint(Var(windowWidth), GreaterOrEqual, Const(0), Required),
		NewConstraint(Var(box1.left), Equal, Const(0), Required),
		NewConstraint(Var(box2.right), Equal, Var(windowWidth), Required),
		NewConstraint(Var(box2.left), GreaterOrEqual, Var(box1.right), Required),
		NewConstraint(Var(box1.left), LessOrEqual, Var(box1.right), Required),
		NewConstraint(Var(box2.left), LessOrEqual, Var(box2.right), Required),
		NewConstraint(Var(box1.right).MinusExpression(Var(box1.left)), Equal, Const(50), Weak),
		NewConstraint(Var(box2.right).MinusExpression(Var(box2.left)), Equal, Const(100), Weak),
	); err != nil {
		t.Fatalf("AddConstraints error = %v, want nil", err)
	}

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

	if err := solver.AddConstraint(NewConstraint(
		Var(box1.right).MinusExpression(Var(box1.left)).Div(50),
		Equal,
		Var(box2.right).MinusExpression(Var(box2.left)).Div(100),
		Medium,
	)); err != nil {
		t.Fatalf("AddConstraint(ratio) error = %v, want nil", err)
	}
	updateValues(solver.FetchChanges())

	assertNear("box1.right after ratio constraint", valueOf(box1.right), 25)
	assertNear("box2.left after ratio constraint", valueOf(box2.left), 25)
	assertNear("windowWidth after ratio constraint", valueOf(windowWidth), 75)
	assertNear("box1.left after ratio constraint", valueOf(box1.left), 0)
	assertNear("box2.right after ratio constraint", valueOf(box2.right), 75)
}
