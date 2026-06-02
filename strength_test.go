package casow

import "testing"

func TestNewStrengthClampsToLegalRange(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected Strength
	}{
		{name: "negative", value: -1, expected: Zero},
		{name: "zero", value: 0, expected: Zero},
		{name: "weak", value: 1, expected: Weak},
		{name: "medium", value: 1_000, expected: Medium},
		{name: "strong", value: 1_000_000, expected: Strong},
		{name: "required", value: 1_001_001_000, expected: Required},
		{name: "above max", value: 1_001_001_001, expected: Required},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStrength(tt.value); got != tt.expected {
				t.Fatalf("NewStrength(%v) = %v, want %v", tt.value, got, tt.expected)
			}
		})
	}
}

func TestCreateStrengthBuildsWeightedStrength(t *testing.T) {
	tests := []struct {
		name       string
		strong     float64
		medium     float64
		weak       float64
		multiplier float64
		expected   Strength
	}{
		{name: "all zero", strong: 0, medium: 0, weak: 0, multiplier: 1, expected: Zero},
		{name: "weak", strong: 0, medium: 0, weak: 1, multiplier: 1, expected: Weak},
		{name: "medium", strong: 0, medium: 1, weak: 0, multiplier: 1, expected: Medium},
		{name: "strong", strong: 1, medium: 0, weak: 0, multiplier: 1, expected: Strong},
		{name: "weak clip", strong: 0, medium: 0, weak: 1000, multiplier: 2, expected: Medium},
		{name: "medium clip", strong: 0, medium: 1000, weak: 0, multiplier: 2, expected: Strong},
		{name: "strong clip", strong: 1000, medium: 0, weak: 0, multiplier: 2, expected: Strong.Mul(1000)},
		{name: "all non-zero", strong: 1, medium: 1, weak: 1, multiplier: 1, expected: Strong.Add(Medium).Add(Weak)},
		{name: "multiplier", strong: 1, medium: 1, weak: 1, multiplier: 2, expected: Strong.Add(Medium).Add(Weak).Mul(2)},
		{name: "max", strong: 1000, medium: 1000, weak: 1000, multiplier: 1, expected: Required},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateStrength(tt.strong, tt.medium, tt.weak, tt.multiplier)
			if got != tt.expected {
				t.Fatalf("CreateStrength(%v, %v, %v, %v) = %v, want %v", tt.strong, tt.medium, tt.weak, tt.multiplier, got, tt.expected)
			}
		})
	}
}

func TestStrengthAdd_shouldMatchUpstreamTableAndClampToLegalRange(t *testing.T) {
	tests := []struct {
		name     string
		lhs      Strength
		rhs      Strength
		expected Strength
	}{
		{name: "zero plus zero", lhs: Zero, rhs: Zero, expected: Zero},
		{name: "zero plus weak", lhs: Zero, rhs: Weak, expected: Weak},
		{name: "weak plus zero", lhs: Weak, rhs: Zero, expected: Weak},
		{name: "weak plus weak", lhs: Weak, rhs: Weak, expected: NewStrength(2.0)},
		{name: "weak plus medium", lhs: Weak, rhs: Medium, expected: NewStrength(1001.0)},
		{name: "medium plus strong", lhs: Medium, rhs: Strong, expected: NewStrength(1_001_000.0)},
		{name: "strong plus required", lhs: Strong, rhs: Required, expected: Required},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lhs.Add(tt.rhs); got != tt.expected {
				t.Fatalf("%v.Add(%v) = %v, want %v", tt.lhs, tt.rhs, got, tt.expected)
			}
		})
	}
}

func TestStrengthSub_shouldMatchUpstreamTableAndClampToLegalRange(t *testing.T) {
	tests := []struct {
		name     string
		lhs      Strength
		rhs      Strength
		expected Strength
	}{
		{name: "saturate low", lhs: Zero, rhs: Weak, expected: Zero},
		{name: "zero minus zero", lhs: Zero, rhs: Zero, expected: Zero},
		{name: "weak minus zero", lhs: Weak, rhs: Zero, expected: Weak},
		{name: "weak minus weak", lhs: Weak, rhs: Weak, expected: Zero},
		{name: "medium minus weak", lhs: Medium, rhs: Weak, expected: NewStrength(999.0)},
		{name: "strong minus medium", lhs: Strong, rhs: Medium, expected: NewStrength(999_000.0)},
		{name: "required minus strong", lhs: Required, rhs: Strong, expected: NewStrength(1_000_001_000.0)},
		{name: "required minus required", lhs: Required, rhs: Required, expected: Zero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lhs.Sub(tt.rhs); got != tt.expected {
				t.Fatalf("%v.Sub(%v) = %v, want %v", tt.lhs, tt.rhs, got, tt.expected)
			}
		})
	}
}

func TestStrengthMul_shouldMatchUpstreamTableAndClampToLegalRange(t *testing.T) {
	tests := []struct {
		name       string
		lhs        Strength
		multiplier float64
		expected   Strength
	}{
		{name: "negative", lhs: Weak, multiplier: -1.0, expected: Zero},
		{name: "zero mul zero", lhs: Zero, multiplier: 0.0, expected: Zero},
		{name: "zero mul one", lhs: Zero, multiplier: 1.0, expected: Zero},
		{name: "weak mul zero", lhs: Weak, multiplier: 0.0, expected: Zero},
		{name: "weak mul one", lhs: Weak, multiplier: 1.0, expected: Weak},
		{name: "weak mul two", lhs: Weak, multiplier: 2.0, expected: NewStrength(2.0)},
		{name: "medium mul half", lhs: Medium, multiplier: 0.5, expected: NewStrength(500.0)},
		{name: "strong mul two", lhs: Strong, multiplier: 2.0, expected: NewStrength(2_000_000.0)},
		{name: "required mul half", lhs: Required, multiplier: 0.5, expected: NewStrength(500_500_500.0)},
		{name: "required mul two clamps high", lhs: Required, multiplier: 2.0, expected: Required},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lhs.Mul(tt.multiplier); got != tt.expected {
				t.Fatalf("%v.Mul(%v) = %v, want %v", tt.lhs, tt.multiplier, got, tt.expected)
			}
		})
	}
}

func TestStrengthDiv_shouldMatchClampSemantics(t *testing.T) {
	tests := []struct {
		name     string
		lhs      Strength
		divisor  float64
		expected Strength
	}{
		{name: "negative", lhs: Weak, divisor: -1.0, expected: Zero},
		{name: "weak div one", lhs: Weak, divisor: 1.0, expected: Weak},
		{name: "weak div half", lhs: Weak, divisor: 0.5, expected: NewStrength(2.0)},
		{name: "medium div two", lhs: Medium, divisor: 2.0, expected: NewStrength(500.0)},
		{name: "strong div two", lhs: Strong, divisor: 2.0, expected: NewStrength(500_000.0)},
		{name: "required div two", lhs: Required, divisor: 2.0, expected: NewStrength(500_500_500.0)},
		{name: "required div half clamps high", lhs: Required, divisor: 0.5, expected: Required},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lhs.Div(tt.divisor); got != tt.expected {
				t.Fatalf("%v.Div(%v) = %v, want %v", tt.lhs, tt.divisor, got, tt.expected)
			}
		})
	}
}

func TestStrengthComparisonIsDeterministic(t *testing.T) {
	tests := []struct {
		name       string
		lhs        Strength
		rhs        Strength
		compare    int
		expectLess bool
	}{
		{name: "less", lhs: Weak, rhs: Medium, compare: -1, expectLess: true},
		{name: "equal", lhs: Medium, rhs: Medium, compare: 0, expectLess: false},
		{name: "greater", lhs: Strong, rhs: Medium, compare: 1, expectLess: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lhs.Compare(tt.rhs); got != tt.compare {
				t.Fatalf("%v.Compare(%v) = %d, want %d", tt.lhs, tt.rhs, got, tt.compare)
			}
			if got := tt.lhs.Less(tt.rhs); got != tt.expectLess {
				t.Fatalf("%v.Less(%v) = %t, want %t", tt.lhs, tt.rhs, got, tt.expectLess)
			}
		})
	}
}

func TestStrengthValueReturnsRawValue(t *testing.T) {
	if got := Strong.Value(); got != 1_000_000 {
		t.Fatalf("Strong.Value() = %v, want 1000000", got)
	}
}
