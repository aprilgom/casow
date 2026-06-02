package casow

const (
	maxStrength      = 1_001_001_000.0
	componentMax     = 1000.0
	strongMultiplier = 1_000_000.0
	mediumMultiplier = 1_000.0
	weakMultiplier   = 1.0
)

type Strength struct {
	value float64
}

var (
	Required = Strength{value: maxStrength}
	Strong   = Strength{value: strongMultiplier}
	Medium   = Strength{value: mediumMultiplier}
	Weak     = Strength{value: weakMultiplier}
	Zero     = Strength{value: 0}
)

func NewStrength(value float64) Strength {
	return Strength{value: clampStrengthValue(value, 0, maxStrength)}
}

func CreateStrength(strong, medium, weak, multiplier float64) Strength {
	strongValue := clampStrengthValue(strong*multiplier, 0, componentMax) * Strong.value
	mediumValue := clampStrengthValue(medium*multiplier, 0, componentMax) * Medium.value
	weakValue := clampStrengthValue(weak*multiplier, 0, componentMax) * Weak.value

	return NewStrength(strongValue + mediumValue + weakValue)
}

func (s Strength) Value() float64 {
	return s.value
}

func (s Strength) Add(other Strength) Strength {
	return NewStrength(s.value + other.value)
}

func (s Strength) Sub(other Strength) Strength {
	return NewStrength(s.value - other.value)
}

func (s Strength) Mul(multiplier float64) Strength {
	return NewStrength(s.value * multiplier)
}

func (s Strength) Div(divisor float64) Strength {
	return NewStrength(s.value / divisor)
}

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
