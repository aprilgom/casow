package casow

import (
	"errors"
	"testing"
)

func TestPublicErrorSentinels(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "duplicate constraint",
			err:  ErrDuplicateConstraint,
			want: "constraint has already been added to the solver",
		},
		{
			name: "unsatisfiable constraint",
			err:  ErrUnsatisfiableConstraint,
			want: "required constraint is unsatisfiable with the existing constraints",
		},
		{
			name: "unknown constraint",
			err:  ErrUnknownConstraint,
			want: "constraint was not already in the solver",
		},
		{
			name: "duplicate edit variable",
			err:  ErrDuplicateEditVariable,
			want: "variable is already marked as an edit variable",
		},
		{
			name: "unknown edit variable",
			err:  ErrUnknownEditVariable,
			want: "variable was not an edit variable",
		},
		{
			name: "bad required strength",
			err:  ErrBadRequiredStrength,
			want: "required strength is illegal for edit variables",
		},
		{
			name: "internal solver error",
			err:  ErrInternalSolver,
			want: "solver entered an invalid internal state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatalf("sentinel error is nil")
			}

			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("Error() = %q, want %q", got, tt.want)
			}

			if !errors.Is(tt.err, tt.err) {
				t.Fatalf("errors.Is(%v, %v) = false, want true", tt.err, tt.err)
			}
		})
	}
}
