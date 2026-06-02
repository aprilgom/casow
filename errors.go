package casow

import "errors"

var (
	// ErrDuplicateConstraint reports that a constraint was already added.
	ErrDuplicateConstraint = errors.New("constraint has already been added to the solver")
	// ErrUnsatisfiableConstraint reports that required constraints conflict.
	ErrUnsatisfiableConstraint = errors.New("required constraint is unsatisfiable with the existing constraints")
	// ErrUnknownConstraint reports that a constraint is not in the solver.
	ErrUnknownConstraint = errors.New("constraint was not already in the solver")
	// ErrDuplicateEditVariable reports that an edit variable already exists.
	ErrDuplicateEditVariable = errors.New("variable is already marked as an edit variable")
	// ErrUnknownEditVariable reports that a variable is not an edit variable.
	ErrUnknownEditVariable = errors.New("variable was not an edit variable")
	// ErrBadRequiredStrength reports that Required was used for an edit variable.
	ErrBadRequiredStrength = errors.New("required strength is illegal for edit variables")
	// ErrInternalSolver reports an invalid internal solver state.
	ErrInternalSolver = errors.New("solver entered an invalid internal state")
)
