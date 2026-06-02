package casow

import "testing"

func TestNewSymbol_shouldPreserveIDAndKind_whenGivenValues(t *testing.T) {
	s := newSymbol(42, symbolSlack)

	if s.id != 42 {
		t.Fatalf("newSymbol id = %d, want %d", s.id, 42)
	}
	if s.Kind() != symbolSlack {
		t.Fatalf("newSymbol kind = %v, want %v", s.Kind(), symbolSlack)
	}
}

func TestInvalidSymbol_shouldReturnInvalidZeroSymbol(t *testing.T) {
	s := invalidSymbol()

	if s.id != 0 {
		t.Fatalf("invalidSymbol id = %d, want %d", s.id, 0)
	}
	if s.Kind() != symbolInvalid {
		t.Fatalf("invalidSymbol kind = %v, want %v", s.Kind(), symbolInvalid)
	}
}

func TestSymbol_shouldCompareAndWorkAsMapKey_whenCopied(t *testing.T) {
	s := newSymbol(7, symbolError)
	copy := s
	values := map[symbol]string{s: "marker"}

	if copy != s {
		t.Fatalf("copied symbol should compare equal: got %v, want %v", copy, s)
	}
	if got := values[copy]; got != "marker" {
		t.Fatalf("symbol map lookup = %q, want %q", got, "marker")
	}
}

func TestNearZero_shouldMatchRustThreshold(t *testing.T) {
	cases := []struct {
		name  string
		value float64
		want  bool
	}{
		{name: "zero", value: 0, want: true},
		{name: "positive below threshold", value: 0.999999e-8, want: true},
		{name: "negative below threshold", value: -0.999999e-8, want: true},
		{name: "positive threshold", value: 1e-8, want: false},
		{name: "negative threshold", value: -1e-8, want: false},
		{name: "positive above threshold", value: 1.000001e-8, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := nearZero(tc.value); got != tc.want {
				t.Fatalf("nearZero(%g) = %t, want %t", tc.value, got, tc.want)
			}
		})
	}
}

func TestNewRow_shouldStoreConstantAndNoCells_whenGivenConstant(t *testing.T) {
	r := newRow(3.5)

	if r.constant != 3.5 {
		t.Fatalf("newRow constant = %g, want %g", r.constant, 3.5)
	}
	if len(r.cells) != 0 {
		t.Fatalf("newRow cells length = %d, want %d", len(r.cells), 0)
	}
}

func TestAddConstant_shouldAdjustAndReturnConstant(t *testing.T) {
	r := newRow(2.25)

	got := r.AddConstant(-3.5)

	if got != -1.25 {
		t.Fatalf("AddConstant return = %g, want %g", got, -1.25)
	}
	if r.constant != -1.25 {
		t.Fatalf("row constant = %g, want %g", r.constant, -1.25)
	}
}

func TestInsertSymbol_shouldAddAccumulateAndRemoveNearZeroCoefficient(t *testing.T) {
	s := newSymbol(1, symbolSlack)
	r := newRow(0)

	r.InsertSymbol(s, 2.5)
	r.InsertSymbol(s, -1.0)
	r.InsertSymbol(s, -1.5+0.5e-8)

	if got := r.CoefficientFor(s); got != 0 {
		t.Fatalf("coefficient after near-zero accumulation = %g, want %g", got, 0.0)
	}
	if len(r.cells) != 0 {
		t.Fatalf("cells length after near-zero removal = %d, want %d", len(r.cells), 0)
	}
}

func TestInsertSymbol_shouldIgnoreNewNearZeroCoefficient(t *testing.T) {
	s := newSymbol(1, symbolSlack)
	r := newRow(0)

	r.InsertSymbol(s, 0.5e-8)

	if len(r.cells) != 0 {
		t.Fatalf("near-zero insert cells length = %d, want %d", len(r.cells), 0)
	}
}

func TestInsertRow_shouldScaleOtherRowAndReportConstantChange(t *testing.T) {
	x := newSymbol(1, symbolExternal)
	y := newSymbol(2, symbolSlack)
	r := newRow(1)
	r.InsertSymbol(x, 2)
	other := newRow(3)
	other.InsertSymbol(x, -4)
	other.InsertSymbol(y, 5)

	changed := r.InsertRow(&other, 2)

	if !changed {
		t.Fatalf("InsertRow changed = false, want true")
	}
	assertRow(t, r, 7, map[symbol]float64{x: -6, y: 10})
}

func TestInsertRow_shouldReportNoConstantChange_whenScaledConstantIsZero(t *testing.T) {
	x := newSymbol(1, symbolExternal)
	r := newRow(1)
	other := newRow(0)
	other.InsertSymbol(x, 5)

	changed := r.InsertRow(&other, 2)

	if changed {
		t.Fatalf("InsertRow changed = true, want false")
	}
	assertRow(t, r, 1, map[symbol]float64{x: 10})
}

func TestRemoveSymbol_shouldRemoveCell(t *testing.T) {
	s := newSymbol(1, symbolSlack)
	r := newRow(0)
	r.InsertSymbol(s, 4)

	r.RemoveSymbol(s)

	if got := r.CoefficientFor(s); got != 0 {
		t.Fatalf("removed coefficient = %g, want %g", got, 0.0)
	}
}

func TestReverseSign_shouldNegateConstantAndCoefficients(t *testing.T) {
	x := newSymbol(1, symbolExternal)
	y := newSymbol(2, symbolSlack)
	r := newRow(3)
	r.InsertSymbol(x, 4)
	r.InsertSymbol(y, -5)

	r.ReverseSign()

	assertRow(t, r, -3, map[symbol]float64{x: -4, y: 5})
}

func TestSolveForSymbol_shouldRemoveSubjectAndScaleRow(t *testing.T) {
	x := newSymbol(1, symbolExternal)
	y := newSymbol(2, symbolSlack)
	r := newRow(6)
	r.InsertSymbol(x, -2)
	r.InsertSymbol(y, 4)

	r.SolveForSymbol(x)

	assertRow(t, r, 3, map[symbol]float64{y: 2})
}

func TestSolveForSymbols_shouldInsertLhsAndSolveForRhs(t *testing.T) {
	lhs := newSymbol(1, symbolExternal)
	rhs := newSymbol(2, symbolSlack)
	y := newSymbol(3, symbolError)
	r := newRow(6)
	r.InsertSymbol(rhs, -2)
	r.InsertSymbol(y, 4)

	r.SolveForSymbols(lhs, rhs)

	assertRow(t, r, 3, map[symbol]float64{lhs: -0.5, y: 2})
}

func TestCoefficientFor_shouldReturnZero_whenSymbolIsMissing(t *testing.T) {
	r := newRow(0)

	if got := r.CoefficientFor(newSymbol(1, symbolExternal)); got != 0 {
		t.Fatalf("missing coefficient = %g, want %g", got, 0.0)
	}
}

func TestSubstitute_shouldReplaceSymbolWithScaledRowAndReportSubstitution(t *testing.T) {
	x := newSymbol(1, symbolExternal)
	y := newSymbol(2, symbolSlack)
	z := newSymbol(3, symbolError)
	r := newRow(1)
	r.InsertSymbol(x, 2)
	r.InsertSymbol(y, 3)
	replacement := newRow(4)
	replacement.InsertSymbol(y, -1)
	replacement.InsertSymbol(z, 5)

	substituted := r.Substitute(x, &replacement)

	if !substituted {
		t.Fatalf("Substitute returned false, want true")
	}
	assertRow(t, r, 9, map[symbol]float64{y: 1, z: 10})
}

func TestSubstitute_shouldReturnFalse_whenSymbolIsMissing(t *testing.T) {
	x := newSymbol(1, symbolExternal)
	r := newRow(1)
	replacement := newRow(4)

	substituted := r.Substitute(x, &replacement)

	if substituted {
		t.Fatalf("Substitute returned true, want false")
	}
	assertRow(t, r, 1, nil)
}

func assertRow(t *testing.T, r row, wantConstant float64, wantCells map[symbol]float64) {
	t.Helper()
	if r.constant != wantConstant {
		t.Fatalf("row constant = %g, want %g", r.constant, wantConstant)
	}
	if len(r.cells) != len(wantCells) {
		t.Fatalf("row cells length = %d, want %d; cells=%v", len(r.cells), len(wantCells), r.cells)
	}
	for s, want := range wantCells {
		if got := r.CoefficientFor(s); got != want {
			t.Fatalf("coefficient for %v = %g, want %g", s, got, want)
		}
	}
}
