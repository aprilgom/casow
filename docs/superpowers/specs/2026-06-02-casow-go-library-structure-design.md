# Casow Go Library Structure Design

## Purpose

Casow is the Go port of the Rust `kasuari` Cassowary constraint solver. The first milestone is a Go library package, not a CLI application. The Rust project remains the behavioral reference for algorithm details, public concepts, and test expectations.

## Scope

This design covers the initial repository and package structure for the Go port. It does not implement the solver yet. The next implementation plan should create the package skeleton, public API types, and the first focused tests.

## Recommended Approach

Use a single root package named `casow`.

Go users should be able to import one package and access the core solver concepts directly:

```go
import "github.com/aprilgom/casow"
```

This keeps the early port easy to compare against `kasuari` while avoiding unnecessary package fragmentation. Cassowary concepts such as variables, expressions, terms, constraints, strengths, relations, and the solver are tightly coupled, so splitting them into separate public packages would add import friction without improving the first milestone.

## Initial File Structure

```text
casow/
  go.mod
  README.md
  AGENTS.md
  variable.go
  term.go
  expression.go
  constraint.go
  relation.go
  strength.go
  solver.go
  row.go
  symbol.go
  errors.go
  expression_test.go
  constraint_test.go
  solver_test.go
  cmd/casow/main.go
```

## Package Boundaries

Public API files:

- `variable.go`: `Variable` identity and construction.
- `term.go`: `Term` as a variable/coefficient pair.
- `expression.go`: linear expression construction and arithmetic helpers.
- `constraint.go`: constraint construction and immutable constraint data.
- `relation.go`: relational operators such as equality, less-than-or-equal, and greater-than-or-equal.
- `strength.go`: Cassowary strengths and validation.
- `solver.go`: public solver operations.
- `errors.go`: typed public errors.

Private implementation files:

- `row.go`: tableau row operations.
- `symbol.go`: internal symbol identifiers and symbol kinds.

The executable at `cmd/casow/main.go` should remain minimal until there is a concrete CLI use case. It should not drive the library design.

## Public API Direction

Use Go-idiomatic explicit constructors and methods. Go cannot reproduce Rust's operator-overload syntax, so the API should be clear rather than syntactically similar.

Target usage shape:

```go
x := casow.NewVariable()
y := casow.NewVariable()

solver := casow.NewSolver()
err := solver.AddConstraint(
    casow.NewConstraint(
        casow.Expr(x),
        casow.EQ,
        casow.Expr(y).PlusConstant(10),
        casow.Required,
    ),
)
```

The exact method names can be refined during implementation, but the API should stay explicit, readable, and easy to test.

## Testing Strategy

Port tests before porting implementation code. Each implementation slice should follow this order:

1. Translate the relevant Rust test or Rust unit-test expectation into a Go `*_test.go` test.
2. Run the Go test and verify that it fails for the expected missing behavior.
3. Port only the minimum Go code needed for that test.
4. Run the targeted test and then `go test ./...`.
5. Refactor only after the tests are green.

Start with small unit tests for value types before porting full solver behavior:

1. `Strength` constants and comparison behavior.
2. `Variable`, `Term`, and `Expression` construction and arithmetic helpers.
3. `Constraint` construction and relation handling.
4. Solver smoke tests from `kasuari/tests`.
5. Removal/edit-variable behavior from `kasuari/tests/removal.rs` and related Rust tests.

Tests should use Go's standard `testing` package and live next to the package files.

## Error Handling

Model public errors as named Go errors or small typed error values that map to kasuari's public error categories:

- add constraint failures
- remove constraint failures
- add edit variable failures
- remove edit variable failures
- suggest value failures
- internal solver failures

The implementation plan should choose concrete Go error shapes once the solver method signatures are defined.

## Non-Goals

- No multi-package public API for the first milestone.
- No CLI behavior beyond a minimal command entry point.
- No UI/layout wrapper abstractions.
- No dependency on the Rust crate at runtime.

## Approval

The approved direction is option 1: a single root Go library package named `casow` with Go-idiomatic explicit APIs.
