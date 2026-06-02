package casow

import "sync/atomic"

var nextConstraintID atomic.Uint64

type Constraint struct {
	id         uint64
	expression *Expression
	operator   RelationalOperator
	strength   Strength
}

func NewConstraint(lhs Expression, op RelationalOperator, rhs Expression, strength Strength) Constraint {
	expression := lhs.MinusExpression(rhs)
	return Constraint{
		id:         nextConstraintID.Add(1) - 1,
		expression: &expression,
		operator:   op,
		strength:   strength,
	}
}

func (c Constraint) Expression() Expression {
	if c.expression == nil {
		return Expression{}
	}
	return *c.expression
}

func (c Constraint) Operator() RelationalOperator {
	return c.operator
}

func (c Constraint) Strength() Strength {
	return c.strength
}
