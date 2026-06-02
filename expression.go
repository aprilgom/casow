package casow

type Expression struct {
	terms    []Term
	constant float64
}

func NewExpression(terms []Term, constant float64) Expression {
	return Expression{terms: copyTerms(terms), constant: constant}
}

func ConstantExpression(constant float64) Expression {
	return NewExpression(nil, constant)
}

func ExpressionFromTerm(term Term) Expression {
	return NewExpression([]Term{term}, 0.0)
}

func ExpressionFromTerms(terms []Term) Expression {
	return NewExpression(terms, 0.0)
}

func ExpressionFromVariable(variable Variable) Expression {
	return ExpressionFromTerm(TermFromVariable(variable))
}

func (e Expression) Terms() []Term {
	return copyTerms(e.terms)
}

func (e Expression) Constant() float64 {
	return e.constant
}

func (e Expression) Negate() Expression {
	return e.Mul(-1.0)
}

func (e Expression) Mul(multiplier float64) Expression {
	terms := make([]Term, len(e.terms))
	for i, term := range e.terms {
		terms[i] = term.Mul(multiplier)
	}
	return NewExpression(terms, e.constant*multiplier)
}

func (e Expression) Div(divisor float64) Expression {
	terms := make([]Term, len(e.terms))
	for i, term := range e.terms {
		terms[i] = term.Div(divisor)
	}
	return NewExpression(terms, e.constant/divisor)
}

func (e Expression) PlusConstant(value float64) Expression {
	return NewExpression(e.terms, e.constant+value)
}

func (e Expression) MinusConstant(value float64) Expression {
	return NewExpression(e.terms, e.constant-value)
}

func (e Expression) PlusExpression(other Expression) Expression {
	terms := make([]Term, 0, len(e.terms)+len(other.terms))
	terms = append(terms, e.terms...)
	terms = append(terms, other.terms...)
	return NewExpression(terms, e.constant+other.constant)
}

func (e Expression) MinusExpression(other Expression) Expression {
	terms := make([]Term, 0, len(e.terms)+len(other.terms))
	terms = append(terms, e.terms...)
	for _, term := range other.terms {
		terms = append(terms, term.Negate())
	}
	return NewExpression(terms, e.constant-other.constant)
}

func copyTerms(terms []Term) []Term {
	if len(terms) == 0 {
		return nil
	}
	copied := make([]Term, len(terms))
	copy(copied, terms)
	return copied
}
