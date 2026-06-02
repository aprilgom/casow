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

func TestSolverGetValue_shouldReturnZero_whenVariableIsUnknown(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()

	if got := solver.GetValue(x); got != 0 {
		t.Fatalf("GetValue(x) = %v, want 0", got)
	}
}

func TestSolverGetValue_shouldReturnSolvedValue_whenConstraintAdded(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()

	if err := solver.AddConstraint(NewConstraint(ExpressionFromVariable(x), Equal, ConstantExpression(10), Required)); err != nil {
		t.Fatalf("AddConstraint error = %v, want nil", err)
	}

	if got := solver.GetValue(x); got != 10 {
		t.Fatalf("GetValue(x) = %v, want 10", got)
	}
}

func TestSolverGetValue_shouldReturnSuggestedValue_whenEditVariableSuggested(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()

	if err := solver.AddEditVariable(x, Strong); err != nil {
		t.Fatalf("AddEditVariable error = %v, want nil", err)
	}
	if err := solver.SuggestValue(x, 42); err != nil {
		t.Fatalf("SuggestValue error = %v, want nil", err)
	}

	if got := solver.GetValue(x); got != 42 {
		t.Fatalf("GetValue(x) = %v, want 42", got)
	}
}

func TestSolverReset_shouldClearConstraintsEditsAndValues(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()
	y := NewVariable()
	constraint := NewConstraint(ExpressionFromVariable(x), Equal, ConstantExpression(10), Required)

	if err := solver.AddConstraint(constraint); err != nil {
		t.Fatalf("AddConstraint error = %v, want nil", err)
	}
	if err := solver.AddEditVariable(y, Strong); err != nil {
		t.Fatalf("AddEditVariable error = %v, want nil", err)
	}
	if err := solver.SuggestValue(y, 20); err != nil {
		t.Fatalf("SuggestValue(first) error = %v, want nil", err)
	}

	solver.Reset()

	if got := solver.GetValue(x); got != 0 {
		t.Fatalf("GetValue(x) after Reset = %v, want 0", got)
	}
	if got := solver.GetValue(y); got != 0 {
		t.Fatalf("GetValue(y) after Reset = %v, want 0", got)
	}
	if solver.HasEditVariable(y) {
		t.Fatal("HasEditVariable(y) after Reset = true, want false")
	}
	if err := solver.RemoveConstraint(constraint); err != ErrUnknownConstraint {
		t.Fatalf("RemoveConstraint after Reset error = %v, want %v", err, ErrUnknownConstraint)
	}
	if err := solver.SuggestValue(y, 30); err != ErrUnknownEditVariable {
		t.Fatalf("SuggestValue after Reset error = %v, want %v", err, ErrUnknownEditVariable)
	}

	if err := solver.AddConstraint(NewConstraint(ExpressionFromVariable(x), Equal, ConstantExpression(15), Required)); err != nil {
		t.Fatalf("AddConstraint after Reset error = %v, want nil", err)
	}
	if got := solver.GetValue(x); got != 15 {
		t.Fatalf("GetValue(x) after new constraint = %v, want 15", got)
	}
}

func TestSolverAddConstraint_shouldRejectDuplicateConstraint(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()
	constraint := NewConstraint(ExpressionFromVariable(x), Equal, ConstantExpression(10), Required)

	if err := solver.AddConstraint(constraint); err != nil {
		t.Fatalf("AddConstraint(first) error = %v, want nil", err)
	}
	if err := solver.AddConstraint(constraint); err != ErrDuplicateConstraint {
		t.Fatalf("AddConstraint(second) error = %v, want %v", err, ErrDuplicateConstraint)
	}
}

func TestSolverRemoveConstraint_shouldRejectUnknownConstraint(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()
	constraint := NewConstraint(ExpressionFromVariable(x), Equal, ConstantExpression(10), Required)

	if err := solver.RemoveConstraint(constraint); err != ErrUnknownConstraint {
		t.Fatalf("RemoveConstraint error = %v, want %v", err, ErrUnknownConstraint)
	}
}

func TestSolverAddConstraint_shouldRejectUnsatisfiableRequiredEqualities(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()

	if err := solver.AddConstraint(NewConstraint(ExpressionFromVariable(x), Equal, ConstantExpression(10), Required)); err != nil {
		t.Fatalf("AddConstraint(x == 10) error = %v, want nil", err)
	}
	if err := solver.AddConstraint(NewConstraint(ExpressionFromVariable(x), Equal, ConstantExpression(20), Required)); err != ErrUnsatisfiableConstraint {
		t.Fatalf("AddConstraint(x == 20) error = %v, want %v", err, ErrUnsatisfiableConstraint)
	}
}

func TestSolverAddConstraints_shouldAddAllConstraints_whenAllAreSatisfiable(t *testing.T) {
	valueOf, updateValues := newValues()

	solver := NewSolver()
	x := NewVariable()
	y := NewVariable()

	if err := solver.AddConstraints(
		NewConstraint(ExpressionFromVariable(x), Equal, ConstantExpression(10), Required),
		NewConstraint(ExpressionFromVariable(y), Equal, ConstantExpression(20), Required),
	); err != nil {
		t.Fatalf("AddConstraints error = %v, want nil", err)
	}
	updateValues(solver.FetchChanges())

	if got := valueOf(x); got != 10 {
		t.Fatalf("valueOf(x) = %v, want 10", got)
	}
	if got := valueOf(y); got != 20 {
		t.Fatalf("valueOf(y) = %v, want 20", got)
	}
}

func TestSolverAddConstraints_shouldReturnFirstErrorWithoutRollback_whenConstraintFails(t *testing.T) {
	valueOf, updateValues := newValues()

	solver := NewSolver()
	x := NewVariable()

	if err := solver.AddConstraints(
		NewConstraint(ExpressionFromVariable(x), Equal, ConstantExpression(10), Required),
		NewConstraint(ExpressionFromVariable(x), Equal, ConstantExpression(20), Required),
	); err != ErrUnsatisfiableConstraint {
		t.Fatalf("AddConstraints error = %v, want %v", err, ErrUnsatisfiableConstraint)
	}
	updateValues(solver.FetchChanges())

	if got := valueOf(x); got != 10 {
		t.Fatalf("valueOf(x) after failed batch = %v, want 10", got)
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

func TestSolverRemoveEditVariable_shouldRejectUnknownEditVariable(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()

	if err := solver.RemoveEditVariable(x); err != ErrUnknownEditVariable {
		t.Fatalf("RemoveEditVariable error = %v, want %v", err, ErrUnknownEditVariable)
	}
}

func TestSolverRemoveEditVariable_shouldClearEditVariable_whenEditVariableRemoved(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()

	if err := solver.AddEditVariable(x, Strong); err != nil {
		t.Fatalf("AddEditVariable error = %v, want nil", err)
	}
	if !solver.HasEditVariable(x) {
		t.Fatal("HasEditVariable before remove = false, want true")
	}
	if err := solver.RemoveEditVariable(x); err != nil {
		t.Fatalf("RemoveEditVariable(first) error = %v, want nil", err)
	}
	if solver.HasEditVariable(x) {
		t.Fatal("HasEditVariable after remove = true, want false")
	}
	if err := solver.RemoveEditVariable(x); err != ErrUnknownEditVariable {
		t.Fatalf("RemoveEditVariable(second) error = %v, want %v", err, ErrUnknownEditVariable)
	}
}

func TestSolverSuggestValue_shouldRejectUnknownEditVariableAfterRemove(t *testing.T) {
	solver := NewSolver()
	x := NewVariable()

	if err := solver.AddEditVariable(x, Strong); err != nil {
		t.Fatalf("AddEditVariable error = %v, want nil", err)
	}
	if err := solver.RemoveEditVariable(x); err != nil {
		t.Fatalf("RemoveEditVariable error = %v, want nil", err)
	}
	if err := solver.SuggestValue(x, 42); err != ErrUnknownEditVariable {
		t.Fatalf("SuggestValue after remove error = %v, want %v", err, ErrUnknownEditVariable)
	}
}

func TestSolverQuadrilateral_shouldMatchKasuariReference_whenEditingCorner(t *testing.T) {
	type point struct {
		x Variable
		y Variable
	}

	newPoint := func() point {
		return point{x: NewVariable(), y: NewVariable()}
	}
	variable := func(v Variable) Expression {
		return ExpressionFromVariable(v)
	}
	average := func(a, b Variable) Expression {
		return variable(a).PlusExpression(variable(b)).Div(2)
	}
	addConstraint := func(solver *Solver, lhs Expression, op RelationalOperator, rhs Expression, strength Strength) {
		t.Helper()
		if err := solver.AddConstraint(NewConstraint(lhs, op, rhs, strength)); err != nil {
			t.Fatalf("AddConstraint(%v) error = %v, want nil", op, err)
		}
	}
	assertPoints := func(name string, valueOf func(Variable) float64, points [4]point, want [4][2]float64) {
		t.Helper()
		got := [4][2]float64{
			{valueOf(points[0].x), valueOf(points[0].y)},
			{valueOf(points[1].x), valueOf(points[1].y)},
			{valueOf(points[2].x), valueOf(points[2].y)},
			{valueOf(points[3].x), valueOf(points[3].y)},
		}
		if got != want {
			t.Fatalf("%s = %v, want %v", name, got, want)
		}
	}

	valueOf, updateValues := newValues()
	points := [4]point{newPoint(), newPoint(), newPoint(), newPoint()}
	pointStarts := [4][2]float64{{10, 10}, {10, 200}, {200, 200}, {200, 10}}
	midpoints := [4]point{newPoint(), newPoint(), newPoint(), newPoint()}
	solver := NewSolver()

	weight := 1.0
	for i := range points {
		addConstraint(solver, variable(points[i].x), Equal, ConstantExpression(pointStarts[i][0]), Weak.Mul(weight))
		addConstraint(solver, variable(points[i].y), Equal, ConstantExpression(pointStarts[i][1]), Weak.Mul(weight))
		weight *= 2
	}

	for _, edge := range [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 0}} {
		start, end := edge[0], edge[1]
		addConstraint(solver, variable(midpoints[start].x), Equal, average(points[start].x, points[end].x), Required)
		addConstraint(solver, variable(midpoints[start].y), Equal, average(points[start].y, points[end].y), Required)
	}

	addConstraint(solver, variable(points[0].x).PlusConstant(20), LessOrEqual, variable(points[2].x), Strong)
	addConstraint(solver, variable(points[0].x).PlusConstant(20), LessOrEqual, variable(points[3].x), Strong)
	addConstraint(solver, variable(points[1].x).PlusConstant(20), LessOrEqual, variable(points[2].x), Strong)
	addConstraint(solver, variable(points[1].x).PlusConstant(20), LessOrEqual, variable(points[3].x), Strong)
	addConstraint(solver, variable(points[0].y).PlusConstant(20), LessOrEqual, variable(points[1].y), Strong)
	addConstraint(solver, variable(points[0].y).PlusConstant(20), LessOrEqual, variable(points[2].y), Strong)
	addConstraint(solver, variable(points[3].y).PlusConstant(20), LessOrEqual, variable(points[1].y), Strong)
	addConstraint(solver, variable(points[3].y).PlusConstant(20), LessOrEqual, variable(points[2].y), Strong)

	for _, point := range points {
		addConstraint(solver, variable(point.x), GreaterOrEqual, ConstantExpression(0), Required)
		addConstraint(solver, variable(point.y), GreaterOrEqual, ConstantExpression(0), Required)
		addConstraint(solver, variable(point.x), LessOrEqual, ConstantExpression(500), Required)
		addConstraint(solver, variable(point.y), LessOrEqual, ConstantExpression(500), Required)
	}

	updateValues(solver.FetchChanges())
	assertPoints("initial midpoints", valueOf, midpoints, [4][2]float64{{10, 105}, {105, 200}, {200, 105}, {105, 10}})

	if err := solver.AddEditVariable(points[2].x, Strong); err != nil {
		t.Fatalf("AddEditVariable(points[2].x) error = %v, want nil", err)
	}
	if err := solver.AddEditVariable(points[2].y, Strong); err != nil {
		t.Fatalf("AddEditVariable(points[2].y) error = %v, want nil", err)
	}
	if err := solver.SuggestValue(points[2].x, 300); err != nil {
		t.Fatalf("SuggestValue(points[2].x) error = %v, want nil", err)
	}
	if err := solver.SuggestValue(points[2].y, 400); err != nil {
		t.Fatalf("SuggestValue(points[2].y) error = %v, want nil", err)
	}

	updateValues(solver.FetchChanges())
	assertPoints("edited points", valueOf, points, [4][2]float64{{10, 10}, {10, 200}, {300, 400}, {200, 10}})
	assertPoints("edited midpoints", valueOf, midpoints, [4][2]float64{{10, 105}, {155, 300}, {250, 205}, {105, 10}})
}
