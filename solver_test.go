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
