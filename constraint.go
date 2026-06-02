package casow

import "sync/atomic"

var nextConstraintID atomic.Uint64

// Constraint is a linear relationship enforced by the solver.
//
// Constraints compare a canonical expression to zero. Separate constraints
// with identical expressions are distinct handles so they can be added and
// removed independently.
type Constraint struct {
	id         uint64
	expression *Expression
	operator   RelationalOperator
	strength   Strength
}

// NewConstraint creates a constraint equivalent to lhs op rhs with strength.
func NewConstraint(lhs Expression, op RelationalOperator, rhs Expression, strength Strength) Constraint {
	expression := lhs.MinusExpression(rhs)
	return Constraint{
		id:         nextConstraintID.Add(1) - 1,
		expression: &expression,
		operator:   op,
		strength:   strength,
	}
}

// Expression returns the canonical expression lhs-rhs.
func (c Constraint) Expression() Expression {
	if c.expression == nil {
		return Expression{}
	}
	return *c.expression
}

// Operator returns the constraint relational operator.
func (c Constraint) Operator() RelationalOperator {
	return c.operator
}

// Strength returns the constraint strength.
func (c Constraint) Strength() Strength {
	return c.strength
}
