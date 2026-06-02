package casow

type RelationalOperator int

const (
	LessOrEqual RelationalOperator = iota
	Equal
	GreaterOrEqual
)

func (op RelationalOperator) String() string {
	switch op {
	case LessOrEqual:
		return "<="
	case Equal:
		return "=="
	case GreaterOrEqual:
		return ">="
	default:
		return ""
	}
}

type WeightedRelation struct {
	operator RelationalOperator
	strength Strength
}

func EQ(strength Strength) WeightedRelation {
	return WeightedRelation{operator: Equal, strength: strength}
}

func LE(strength Strength) WeightedRelation {
	return WeightedRelation{operator: LessOrEqual, strength: strength}
}

func GE(strength Strength) WeightedRelation {
	return WeightedRelation{operator: GreaterOrEqual, strength: strength}
}

func (r WeightedRelation) Operator() RelationalOperator {
	return r.operator
}

func (r WeightedRelation) Strength() Strength {
	return r.strength
}
