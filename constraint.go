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
//
// lhs and rhs may be Expression, Variable, Term, float64, float32, or any
// integer type. Unsupported values panic.
func NewConstraint(lhs any, op RelationalOperator, rhs any, strength Strength) Constraint {
	lhsExpression := constraintExpression(lhs)
	rhsExpression := constraintExpression(rhs)
	expression := lhsExpression.MinusExpression(rhsExpression)
	return Constraint{
		id:         nextConstraintID.Add(1),
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

func constraintExpression(value any) Expression {
	switch typed := value.(type) {
	case Expression:
		return typed
	case Variable:
		return ExpressionFromVariable(typed)
	case Term:
		return ExpressionFromTerm(typed)
	case float64:
		return ConstantExpression(typed)
	case float32:
		return ConstantExpression(float64(typed))
	case int:
		return ConstantExpression(float64(typed))
	case int8:
		return ConstantExpression(float64(typed))
	case int16:
		return ConstantExpression(float64(typed))
	case int32:
		return ConstantExpression(float64(typed))
	case int64:
		return ConstantExpression(float64(typed))
	case uint:
		return ConstantExpression(float64(typed))
	case uint8:
		return ConstantExpression(float64(typed))
	case uint16:
		return ConstantExpression(float64(typed))
	case uint32:
		return ConstantExpression(float64(typed))
	case uint64:
		return ConstantExpression(float64(typed))
	default:
		panic("unsupported constraint expression value")
	}
}
