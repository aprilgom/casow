package casow

// Term is a variable multiplied by a coefficient.
type Term struct {
	variable    Variable
	coefficient float64
}

// NewTerm creates a term from a variable and coefficient.
func NewTerm(variable Variable, coefficient float64) Term {
	return Term{variable: variable, coefficient: coefficient}
}

// TermFromVariable creates a unit-coefficient term for variable.
func TermFromVariable(variable Variable) Term {
	return NewTerm(variable, 1.0)
}

// Var returns the variable in the term.
func (t Term) Var() Variable {
	return t.variable
}

// Coefficient returns the coefficient multiplying the term variable.
func (t Term) Coefficient() float64 {
	return t.coefficient
}

// Mul returns the term multiplied by multiplier.
func (t Term) Mul(multiplier float64) Term {
	return NewTerm(t.variable, t.coefficient*multiplier)
}

// Div returns the term divided by divisor.
func (t Term) Div(divisor float64) Term {
	return NewTerm(t.variable, t.coefficient/divisor)
}

// Negate returns the term with its coefficient sign reversed.
func (t Term) Negate() Term {
	return NewTerm(t.variable, -t.coefficient)
}
