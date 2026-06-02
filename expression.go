package casow

// Expression is a linear expression of terms plus a constant.
type Expression struct {
	terms    []Term
	constant float64
}

// NewExpression creates an expression from terms and a constant.
//
// The terms slice is copied so later caller mutations do not affect the
// expression.
func NewExpression(terms []Term, constant float64) Expression {
	return Expression{terms: copyTerms(terms), constant: constant}
}

// ConstantExpression creates an expression containing only a constant.
func ConstantExpression(constant float64) Expression {
	return NewExpression(nil, constant)
}

// Const creates an expression containing only value.
func Const(value float64) Expression {
	return ConstantExpression(value)
}

// ExpressionFromTerm creates an expression from a single term.
func ExpressionFromTerm(term Term) Expression {
	return NewExpression([]Term{term}, 0.0)
}

// ExpressionFromTerms creates an expression from terms with zero constant.
func ExpressionFromTerms(terms []Term) Expression {
	return NewExpression(terms, 0.0)
}

// ExpressionFromVariable creates a unit-coefficient expression for variable.
func ExpressionFromVariable(variable Variable) Expression {
	return ExpressionFromTerm(TermFromVariable(variable))
}

// Var creates a unit-coefficient expression for variable.
func Var(variable Variable) Expression {
	return ExpressionFromVariable(variable)
}

// Terms returns a copy of the expression terms.
func (e Expression) Terms() []Term {
	return copyTerms(e.terms)
}

// Constant returns the expression constant.
func (e Expression) Constant() float64 {
	return e.constant
}

// Negate returns the expression multiplied by -1.
func (e Expression) Negate() Expression {
	return e.Mul(-1.0)
}

// Mul returns the expression multiplied by multiplier.
func (e Expression) Mul(multiplier float64) Expression {
	terms := make([]Term, len(e.terms))
	for i, term := range e.terms {
		terms[i] = term.Mul(multiplier)
	}
	return NewExpression(terms, e.constant*multiplier)
}

// Div returns the expression divided by divisor.
func (e Expression) Div(divisor float64) Expression {
	terms := make([]Term, len(e.terms))
	for i, term := range e.terms {
		terms[i] = term.Div(divisor)
	}
	return NewExpression(terms, e.constant/divisor)
}

// PlusConstant returns the expression plus value.
func (e Expression) PlusConstant(value float64) Expression {
	return NewExpression(e.terms, e.constant+value)
}

// MinusConstant returns the expression minus value.
func (e Expression) MinusConstant(value float64) Expression {
	return NewExpression(e.terms, e.constant-value)
}

// PlusExpression returns the sum of e and other.
func (e Expression) PlusExpression(other Expression) Expression {
	terms := make([]Term, 0, len(e.terms)+len(other.terms))
	terms = append(terms, e.terms...)
	terms = append(terms, other.terms...)
	return NewExpression(terms, e.constant+other.constant)
}

// MinusExpression returns e minus other.
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
