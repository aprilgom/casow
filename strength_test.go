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
		{name: "all non-zero", strong: 1, medium: 1, weak: 1, multiplier: 1, expected: Strong.Add(Medium).Add(Weak)},
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

func TestStrengthArithmeticClampsToLegalRange(t *testing.T) {
	tests := []struct {
		name     string
		got      Strength
		expected Strength
	}{
		{name: "add clamps high", got: Required.Add(Strong), expected: Required},
		{name: "sub clamps low", got: Zero.Sub(Weak), expected: Zero},
		{name: "mul clamps low", got: Weak.Mul(-1), expected: Zero},
		{name: "mul clamps high", got: Required.Mul(2), expected: Required},
		{name: "div clamps low", got: Weak.Div(-1), expected: Zero},
		{name: "div clamps high", got: Required.Div(0.5), expected: Required},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Fatalf("got %v, want %v", tt.got, tt.expected)
			}
		})
	}
}

func TestStrengthArithmeticPreservesIntermediateValues(t *testing.T) {
	tests := []struct {
		name     string
		got      Strength
		expected Strength
	}{
		{name: "add", got: Weak.Add(Medium), expected: NewStrength(1001)},
		{name: "sub", got: Strong.Sub(Medium), expected: NewStrength(999_000)},
		{name: "mul", got: Medium.Mul(0.5), expected: NewStrength(500)},
		{name: "div", got: Strong.Div(2), expected: NewStrength(500_000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Fatalf("got %v, want %v", tt.got, tt.expected)
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
