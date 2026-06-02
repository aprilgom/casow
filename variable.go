package casow

import "sync/atomic"

var nextVariableID atomic.Uint64

// Variable identifies a solver variable.
//
// Variables are immutable value handles. Copying a Variable preserves its
// identity, and each call to NewVariable creates a distinct identity.
type Variable struct {
	id uint64
}

// NewVariable creates a new unique solver variable.
func NewVariable() Variable {
	return Variable{id: nextVariableID.Add(1) - 1}
}

// ID returns the variable's stable process-local identifier.
func (v Variable) ID() uint64 {
	return v.id
}
