package casow

import "testing"

func TestNewVariable_shouldNotEqualZeroValueVariable(t *testing.T) {
	var zero Variable
	nextVariableID.Store(0)

	created := NewVariable()

	if created == zero {
		t.Fatalf("NewVariable() = zero value variable: %v", created)
	}
	if created.ID() == 0 {
		t.Fatalf("NewVariable().ID() = 0, want non-zero")
	}
}

func TestNewVariableReturnsUniqueIdentities(t *testing.T) {
	first := NewVariable()
	second := NewVariable()
	third := NewVariable()

	if first == second || first == third || second == third {
		t.Fatalf("NewVariable returned duplicate identities: %d, %d, %d", first.ID(), second.ID(), third.ID())
	}
}

func TestCopiedVariablePreservesIdentityAndEquality(t *testing.T) {
	original := NewVariable()
	copy := original

	if copy != original {
		t.Fatalf("copied variable should compare equal: got %v, want %v", copy, original)
	}
	if copy.ID() != original.ID() {
		t.Fatalf("copied variable should preserve ID: got %d, want %d", copy.ID(), original.ID())
	}
}

func TestVariableCanBeUsedAsMapKey(t *testing.T) {
	variable := NewVariable()
	values := map[Variable]string{
		variable: "width",
	}

	if got := values[variable]; got != "width" {
		t.Fatalf("map lookup by variable key = %q, want %q", got, "width")
	}
}

func TestNewVariableRepresentsDefaultNewVariableBehavior(t *testing.T) {
	first := NewVariable()
	second := NewVariable()

	if first.ID() == second.ID() {
		t.Fatalf("new variables should receive distinct IDs: got %d and %d", first.ID(), second.ID())
	}
}
