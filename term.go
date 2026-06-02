package casow

type Term struct {
	variable    Variable
	coefficient float64
}

func NewTerm(variable Variable, coefficient float64) Term {
	return Term{variable: variable, coefficient: coefficient}
}

func TermFromVariable(variable Variable) Term {
	return NewTerm(variable, 1.0)
}

func (t Term) Var() Variable {
	return t.variable
}

func (t Term) Coefficient() float64 {
	return t.coefficient
}

func (t Term) Mul(multiplier float64) Term {
	return NewTerm(t.variable, t.coefficient*multiplier)
}

func (t Term) Div(divisor float64) Term {
	return NewTerm(t.variable, t.coefficient/divisor)
}

func (t Term) Negate() Term {
	return NewTerm(t.variable, -t.coefficient)
}
