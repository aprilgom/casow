package casow

import "sync/atomic"

var nextVariableID atomic.Uint64

type Variable struct {
	id uint64
}

func NewVariable() Variable {
	return Variable{id: nextVariableID.Add(1) - 1}
}

func (v Variable) ID() uint64 {
	return v.id
}
