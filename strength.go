package casow

const (
	maxStrength      = 1_001_001_000.0
	componentMax     = 1000.0
	strongMultiplier = 1_000_000.0
	mediumMultiplier = 1_000.0
	weakMultiplier   = 1.0
)

// Strength is a constraint priority.
//
// Larger strengths are preferred over smaller strengths. Values are clamped to
// the legal Cassowary range from Zero to Required.
type Strength struct {
	value float64
}

var (
	// Required is the strongest legal strength and must be satisfied.
	Required = Strength{value: maxStrength}
	// Strong is a high non-required strength.
	Strong = Strength{value: strongMultiplier}
	// Medium is weaker than Strong and stronger than Weak.
	Medium = Strength{value: mediumMultiplier}
	// Weak is the default low preference strength.
	Weak = Strength{value: weakMultiplier}
	// Zero is the weakest legal strength.
	Zero = Strength{value: 0}
)

// NewStrength creates a strength, clamped to the legal range.
func NewStrength(value float64) Strength {
	return Strength{value: clampStrengthValue(value, 0, maxStrength)}
}

// CreateStrength combines strong, medium, and weak components.
func CreateStrength(strong, medium, weak, multiplier float64) Strength {
	strongValue := clampStrengthValue(strong*multiplier, 0, componentMax) * Strong.value
	mediumValue := clampStrengthValue(medium*multiplier, 0, componentMax) * Medium.value
	weakValue := clampStrengthValue(weak*multiplier, 0, componentMax) * Weak.value

	return NewStrength(strongValue + mediumValue + weakValue)
}

// Value returns the numeric strength value.
func (s Strength) Value() float64 {
	return s.value
}

// Add returns s plus other, clamped to the legal range.
func (s Strength) Add(other Strength) Strength {
	return NewStrength(s.value + other.value)
}

// Sub returns s minus other, clamped to the legal range.
func (s Strength) Sub(other Strength) Strength {
	return NewStrength(s.value - other.value)
}

// Mul returns s multiplied by multiplier, clamped to the legal range.
func (s Strength) Mul(multiplier float64) Strength {
	return NewStrength(s.value * multiplier)
}

// Div returns s divided by divisor, clamped to the legal range.
func (s Strength) Div(divisor float64) Strength {
	return NewStrength(s.value / divisor)
}

// Compare compares s with other.
func (s Strength) Compare(other Strength) int {
	switch {
	case s.value < other.value:
		return -1
	case s.value > other.value:
		return 1
	default:
		return 0
	}
}

// Less reports whether s is weaker than other.
func (s Strength) Less(other Strength) bool {
	return s.Compare(other) < 0
}

func clampStrengthValue(value, minValue, maxValue float64) float64 {
	switch {
	case value < minValue:
		return minValue
	case value > maxValue:
		return maxValue
	default:
		return value
	}
}
