package casow

func nearZero(value float64) bool {
	const epsilon = 1e-8
	if value < 0 {
		return -value < epsilon
	}
	return value < epsilon
}

type row struct {
	cells    map[symbol]float64
	constant float64
}

func newRow(constant float64) row {
	return row{
		cells:    make(map[symbol]float64),
		constant: constant,
	}
}

func (r *row) AddConstant(value float64) float64 {
	r.constant += value
	return r.constant
}

func (r *row) InsertSymbol(s symbol, coefficient float64) {
	if current, ok := r.cells[s]; ok {
		next := current + coefficient
		if nearZero(next) {
			delete(r.cells, s)
			return
		}
		r.cells[s] = next
		return
	}

	if !nearZero(coefficient) {
		r.cells[s] = coefficient
	}
}

func (r *row) InsertRow(other *row, coefficient float64) bool {
	constantDiff := other.constant * coefficient
	r.constant += constantDiff
	for s, value := range other.cells {
		r.InsertSymbol(s, value*coefficient)
	}
	return constantDiff != 0
}

func (r *row) RemoveSymbol(s symbol) {
	delete(r.cells, s)
}

func (r *row) ReverseSign() {
	r.constant = -r.constant
	for s, value := range r.cells {
		r.cells[s] = -value
	}
}

func (r *row) SolveForSymbol(s symbol) {
	coefficient := r.cells[s]
	delete(r.cells, s)
	scale := -1.0 / coefficient
	r.constant *= scale
	for s, value := range r.cells {
		r.cells[s] = value * scale
	}
}

func (r *row) SolveForSymbols(lhs, rhs symbol) {
	r.InsertSymbol(lhs, -1.0)
	r.SolveForSymbol(rhs)
}

func (r row) CoefficientFor(s symbol) float64 {
	return r.cells[s]
}

func (r *row) Substitute(s symbol, other *row) bool {
	coefficient, ok := r.cells[s]
	if !ok {
		return false
	}
	delete(r.cells, s)
	return r.InsertRow(other, coefficient)
}
