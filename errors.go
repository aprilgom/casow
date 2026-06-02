package casow

import "errors"

var (
	ErrDuplicateConstraint     = errors.New("constraint has already been added to the solver")
	ErrUnsatisfiableConstraint = errors.New("required constraint is unsatisfiable with the existing constraints")
	ErrUnknownConstraint       = errors.New("constraint was not already in the solver")
	ErrDuplicateEditVariable   = errors.New("variable is already marked as an edit variable")
	ErrUnknownEditVariable     = errors.New("variable was not an edit variable")
	ErrBadRequiredStrength     = errors.New("required strength is illegal for edit variables")
	ErrInternalSolver          = errors.New("solver entered an invalid internal state")
)
