# casow

`casow` is an early Go port of [kasuari](https://github.com/ratatui/kasuari), a
Cassowary constraint solver. Cassowary solves linear constraints incrementally,
which makes it useful for user interface layout systems where some constraints
are required and others are preferences that may be violated when space is
limited.

## Install

The current module path is:

```sh
go get github.com/aprilgom/casow
```

Import it as:

```go
import "github.com/aprilgom/casow"
```

## Usage

This example lays out two horizontal boxes. The first box is aligned to the left,
the second box is aligned to the right, the boxes must not overlap, and each box
has a weak preferred width. The window width is an edit variable so it can be
changed efficiently. After shrinking the window, a medium-strength ratio
constraint asks the boxes to keep the same relative share as their preferred
widths.

```go
package main

import (
	"fmt"

	"github.com/aprilgom/casow"
)

type element struct {
	left  casow.Variable
	right casow.Variable
}

func main() {
	windowWidth := casow.NewVariable()
	box1 := element{left: casow.NewVariable(), right: casow.NewVariable()}
	box2 := element{left: casow.NewVariable(), right: casow.NewVariable()}
	solver := casow.NewSolver()

	if err := solver.AddConstraints(
		casow.NewConstraint(windowWidth, casow.GreaterOrEqual, 0, casow.Required),
		casow.NewConstraint(box1.left, casow.Equal, 0, casow.Required),
		casow.NewConstraint(box2.right, casow.Equal, windowWidth, casow.Required),
		casow.NewConstraint(box2.left, casow.GreaterOrEqual, box1.right, casow.Required),
		casow.NewConstraint(box1.left, casow.LessOrEqual, box1.right, casow.Required),
		casow.NewConstraint(box2.left, casow.LessOrEqual, box2.right, casow.Required),
		casow.NewConstraint(
			casow.Var(box1.right).MinusExpression(casow.Var(box1.left)),
			casow.Equal,
			50,
			casow.Weak,
		),
		casow.NewConstraint(
			casow.Var(box2.right).MinusExpression(casow.Var(box2.left)),
			casow.Equal,
			100,
			casow.Weak,
		),
	); err != nil {
		panic(err)
	}

	if err := solver.AddEditVariable(windowWidth, casow.Strong); err != nil {
		panic(err)
	}
	if err := solver.SuggestValue(windowWidth, 300); err != nil {
		panic(err)
	}

	for _, change := range solver.FetchChanges() {
		fmt.Printf("variable %d = %.0f\n", change.Variable.ID(), change.Value)
	}

	if err := solver.SuggestValue(windowWidth, 75); err != nil {
		panic(err)
	}

	for _, change := range solver.FetchChanges() {
		fmt.Printf("variable %d = %.0f\n", change.Variable.ID(), change.Value)
	}

	if err := solver.AddConstraint(casow.NewConstraint(
		casow.Var(box1.right).MinusExpression(casow.Var(box1.left)).Div(50),
		casow.Equal,
		casow.Var(box2.right).MinusExpression(casow.Var(box2.left)).Div(100),
		casow.Medium,
	)); err != nil {
		panic(err)
	}

	for _, change := range solver.FetchChanges() {
		fmt.Printf("variable %d = %.0f\n", change.Variable.ID(), change.Value)
	}
}
```

When the suggested width is `300`, known values include `windowWidth = 300`,
`box1.right = 50`, `box2.left = 200`, and `box2.right = 300`. When the suggested
width is `75`, the solver preserves the required constraints: the window remains
nonnegative, `box2.right` equals the window width, the boxes do not overlap, and
both boxes have nonnegative widths. The exact weak preferred-width violation can
vary. Adding the medium ratio constraint then gives the boxes proportional
widths at the constrained size, so known values include `box1.right = 25`,
`box2.left = 25`, and `box2.right = 75`.

## Verification

Run the current test suite with:

```sh
go test ./...
```
