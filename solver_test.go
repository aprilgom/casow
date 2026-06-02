package casow

import "testing"

func TestSolverCanBeSharedAcrossGoroutinesByTypeContract(t *testing.T) {
	done := make(chan struct{})

	go func() {
		_ = NewSolver()
		close(done)
	}()

	<-done
}

func TestSolverRemoveConstraint_shouldAllowReplacementConstraint_whenConstraintRemoved(t *testing.T) {
	valueOf, updateValues := newValues()

	solver := NewSolver()
	val := NewVariable()
	constraint := NewConstraint(ExpressionFromVariable(val), Equal, ConstantExpression(100), Required)

	if err := solver.AddConstraint(constraint); err != nil {
		t.Fatalf("AddConstraint(first) error = %v, want nil", err)
	}
	updateValues(solver.FetchChanges())

	if got := valueOf(val); got != 100 {
		t.Fatalf("valueOf(val) after first constraint = %v, want 100", got)
	}

	if err := solver.RemoveConstraint(constraint); err != nil {
		t.Fatalf("RemoveConstraint error = %v, want nil", err)
	}
	if err := solver.AddConstraint(NewConstraint(ExpressionFromVariable(val), Equal, ConstantExpression(0), Required)); err != nil {
		t.Fatalf("AddConstraint(replacement) error = %v, want nil", err)
	}
	updateValues(solver.FetchChanges())

	if got := valueOf(val); got != 0 {
		t.Fatalf("valueOf(val) after replacement constraint = %v, want 0", got)
	}
}

func TestSolverSuggestValue_shouldUpdateEditVariable_whenEditVariableAdded(t *testing.T) {
	valueOf, updateValues := newValues()

	solver := NewSolver()
	x := NewVariable()

	if err := solver.AddEditVariable(x, Strong); err != nil {
		t.Fatalf("AddEditVariable error = %v, want nil", err)
	}
	if err := solver.SuggestValue(x, 42); err != nil {
		t.Fatalf("SuggestValue error = %v, want nil", err)
	}
	updateValues(solver.FetchChanges())

	if got := valueOf(x); got != 42 {
		t.Fatalf("valueOf(x) = %v, want 42", got)
	}
}

func TestSolverAddEditVariable_shouldRejectRequiredStrength(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()

	if err := solver.AddEditVariable(x, Required); err != ErrBadRequiredStrength {
		t.Fatalf("AddEditVariable error = %v, want %v", err, ErrBadRequiredStrength)
	}
}

func TestSolverAddEditVariable_shouldRejectDuplicateEditVariable(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()

	if err := solver.AddEditVariable(x, Strong); err != nil {
		t.Fatalf("AddEditVariable(first) error = %v, want nil", err)
	}
	if err := solver.AddEditVariable(x, Strong); err != ErrDuplicateEditVariable {
		t.Fatalf("AddEditVariable(second) error = %v, want %v", err, ErrDuplicateEditVariable)
	}
}

func TestSolverSuggestValue_shouldRejectUnknownEditVariable(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()

	if err := solver.SuggestValue(x, 42); err != ErrUnknownEditVariable {
		t.Fatalf("SuggestValue error = %v, want %v", err, ErrUnknownEditVariable)
	}
}
