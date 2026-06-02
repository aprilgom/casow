package casow

// RelationalOperator is the comparison operator for a constraint.
type RelationalOperator int

const (
	// LessOrEqual represents <=.
	LessOrEqual RelationalOperator = iota
	// Equal represents ==.
	Equal
	// GreaterOrEqual represents >=.
	GreaterOrEqual
)

// String returns the symbolic operator.
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

// WeightedRelation pairs a relational operator with a strength.
type WeightedRelation struct {
	operator RelationalOperator
	strength Strength
}

// EQ creates an equality relation with strength.
func EQ(strength Strength) WeightedRelation {
	return WeightedRelation{operator: Equal, strength: strength}
}

// LE creates a less-than-or-equal relation with strength.
func LE(strength Strength) WeightedRelation {
	return WeightedRelation{operator: LessOrEqual, strength: strength}
}

// GE creates a greater-than-or-equal relation with strength.
func GE(strength Strength) WeightedRelation {
	return WeightedRelation{operator: GreaterOrEqual, strength: strength}
}

// Operator returns the relation operator.
func (r WeightedRelation) Operator() RelationalOperator {
	return r.operator
}

// Strength returns the relation strength.
func (r WeightedRelation) Strength() Strength {
	return r.strength
}
