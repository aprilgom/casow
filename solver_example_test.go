package casow

func ExampleSolver() {
	type element struct {
		left  Variable
		right Variable
	}

	windowWidth := NewVariable()
	box1 := element{left: NewVariable(), right: NewVariable()}
	box2 := element{left: NewVariable(), right: NewVariable()}

	solver := NewSolver()
	if err := solver.AddConstraints(
		NewConstraint(windowWidth, GreaterOrEqual, 0, Required),
		NewConstraint(box1.left, Equal, 0, Required),
		NewConstraint(box2.right, Equal, windowWidth, Required),
		NewConstraint(box2.left, GreaterOrEqual, box1.right, Required),
		NewConstraint(box1.left, LessOrEqual, box1.right, Required),
		NewConstraint(box2.left, LessOrEqual, box2.right, Required),
		NewConstraint(Var(box1.right).MinusExpression(Var(box1.left)), Equal, 50, Weak),
		NewConstraint(Var(box2.right).MinusExpression(Var(box2.left)), Equal, 100, Weak),
	); err != nil {
		panic(err)
	}

	if err := solver.AddEditVariable(windowWidth, Strong); err != nil {
		panic(err)
	}
	if err := solver.SuggestValue(windowWidth, 300); err != nil {
		panic(err)
	}

	values := make(map[Variable]float64)
	for _, change := range solver.FetchChanges() {
		values[change.Variable] = change.Value
	}

	if err := solver.SuggestValue(windowWidth, 75); err != nil {
		panic(err)
	}
	_ = solver.FetchChanges()

	if err := solver.AddConstraint(NewConstraint(
		Var(box1.right).MinusExpression(Var(box1.left)).Div(50),
		Equal,
		Var(box2.right).MinusExpression(Var(box2.left)).Div(100),
		Medium,
	)); err != nil {
		panic(err)
	}

	for _, change := range solver.FetchChanges() {
		values[change.Variable] = change.Value
	}

	_, _ = values[box1.right], values[box2.left]
}
