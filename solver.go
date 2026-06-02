package casow

import "maps"

import "math"

type tag struct {
	marker symbol
	other  symbol
}

type varData struct {
	value  float64
	symbol symbol
	count  int
}

type editInfo struct {
	tag        tag
	constraint Constraint
	constant   float64
}

type optimizeTarget int

const (
	optimizeObjective optimizeTarget = iota
	optimizeArtificial
)

type Solver struct {
	constraints        map[Constraint]tag
	edits              map[Variable]editInfo
	varData            map[Variable]varData
	varForSymbol       map[symbol]Variable
	publicChanges      []Change
	changed            map[Variable]struct{}
	shouldClearChanges bool
	rows               map[symbol]row
	infeasibleRows     []symbol
	objective          row
	artificial         *row
	idTick             uint64
}

type solverStateSnapshot struct {
	constraints        map[Constraint]tag
	edits              map[Variable]editInfo
	varData            map[Variable]varData
	varForSymbol       map[symbol]Variable
	publicChanges      []Change
	changed            map[Variable]struct{}
	shouldClearChanges bool
	rows               map[symbol]row
	infeasibleRows     []symbol
	objective          row
	artificial         *row
	idTick             uint64
}

// Change records a variable value changed by the solver.
type Change struct {
	Variable Variable
	Value    float64
}

// NewSolver creates an empty constraint solver.
func NewSolver() *Solver {
	return &Solver{
		constraints:  make(map[Constraint]tag),
		edits:        make(map[Variable]editInfo),
		varData:      make(map[Variable]varData),
		varForSymbol: make(map[symbol]Variable),
		changed:      make(map[Variable]struct{}),
		rows:         make(map[symbol]row),
		objective:    newRow(0),
		idTick:       1,
	}
}

// AddConstraint adds a constraint to the solver.
func (s *Solver) AddConstraint(constraint Constraint) error {
	if _, ok := s.constraints[constraint]; ok {
		return ErrDuplicateConstraint
	}
	snapshot := s.snapshotState()
	fail := func(err error) error {
		s.restoreState(snapshot)
		return err
	}

	r, constraintTag := s.createRow(constraint)
	subject := chooseSubject(&r, constraintTag)

	if subject.Kind() == symbolInvalid && allDummies(&r) {
		if !nearZero(r.constant) {
			return fail(ErrUnsatisfiableConstraint)
		}
		subject = constraintTag.marker
	}

	if subject.Kind() == symbolInvalid {
		satisfiable, err := s.addWithArtificialVariable(&r)
		if err != nil {
			return fail(ErrInternalSolver)
		}
		if !satisfiable {
			return fail(ErrUnsatisfiableConstraint)
		}
	} else {
		r.SolveForSymbol(subject)
		s.substitute(subject, &r)
		if subject.Kind() == symbolExternal && r.constant != 0 {
			s.varChanged(s.varForSymbol[subject])
		}
		s.rows[subject] = r
	}

	s.constraints[constraint] = constraintTag
	if err := s.optimize(optimizeObjective); err != nil {
		return fail(ErrInternalSolver)
	}
	return nil
}

// AddConstraints adds constraints sequentially.
//
// If a constraint fails, the error is returned and constraints added earlier in
// the call remain in the solver.
func (s *Solver) AddConstraints(constraints ...Constraint) error {
	for _, constraint := range constraints {
		if err := s.AddConstraint(constraint); err != nil {
			return err
		}
	}
	return nil
}

// RemoveConstraint removes a previously added constraint.
func (s *Solver) RemoveConstraint(constraint Constraint) error {
	constraintTag, ok := s.constraints[constraint]
	if !ok {
		return ErrUnknownConstraint
	}
	delete(s.constraints, constraint)

	s.removeConstraintEffects(constraint, constraintTag)

	if _, ok := s.rows[constraintTag.marker]; ok {
		delete(s.rows, constraintTag.marker)
	} else {
		leaving, r, ok := s.getMarkerLeavingRow(constraintTag.marker)
		if !ok {
			return ErrInternalSolver
		}
		r.SolveForSymbols(leaving, constraintTag.marker)
		s.substitute(constraintTag.marker, &r)
	}

	if err := s.optimize(optimizeObjective); err != nil {
		return ErrInternalSolver
	}

	for _, term := range constraint.Expression().Terms() {
		if nearZero(term.Coefficient()) {
			continue
		}
		v := term.Var()
		data, ok := s.varData[v]
		if !ok {
			continue
		}
		data.count--
		if data.count == 0 {
			delete(s.varForSymbol, data.symbol)
			delete(s.varData, v)
		} else {
			s.varData[v] = data
		}
	}
	return nil
}

// AddEditVariable marks variable as editable with strength.
func (s *Solver) AddEditVariable(variable Variable, strength Strength) error {
	if _, ok := s.edits[variable]; ok {
		return ErrDuplicateEditVariable
	}
	if strength == Required {
		return ErrBadRequiredStrength
	}

	constraint := NewConstraint(ExpressionFromVariable(variable), Equal, ConstantExpression(0), strength)
	if err := s.AddConstraint(constraint); err != nil {
		return err
	}

	s.edits[variable] = editInfo{
		tag:        s.constraints[constraint],
		constraint: constraint,
		constant:   0,
	}
	return nil
}

// RemoveEditVariable removes variable from the edit set.
func (s *Solver) RemoveEditVariable(variable Variable) error {
	info, ok := s.edits[variable]
	if !ok {
		return ErrUnknownEditVariable
	}
	delete(s.edits, variable)

	if err := s.RemoveConstraint(info.constraint); err != nil {
		if err == ErrUnknownConstraint {
			return ErrInternalSolver
		}
		return err
	}
	return nil
}

// HasEditVariable reports whether variable is currently editable.
func (s *Solver) HasEditVariable(variable Variable) bool {
	_, ok := s.edits[variable]
	return ok
}

// GetValue returns the current solved value for variable.
func (s *Solver) GetValue(variable Variable) float64 {
	data, ok := s.varData[variable]
	if !ok {
		return 0
	}
	if r, ok := s.rows[data.symbol]; ok {
		return r.constant
	}
	return 0
}

// Reset clears all solver state.
func (s *Solver) Reset() {
	clear(s.constraints)
	clear(s.edits)
	clear(s.varData)
	clear(s.varForSymbol)
	clear(s.changed)
	clear(s.rows)
	s.publicChanges = s.publicChanges[:0]
	s.infeasibleRows = s.infeasibleRows[:0]
	s.shouldClearChanges = false
	s.objective = newRow(0)
	s.artificial = nil
	s.idTick = 1
}

// SuggestValue suggests a new value for an edit variable.
func (s *Solver) SuggestValue(variable Variable, value float64) error {
	info, ok := s.edits[variable]
	if !ok {
		return ErrUnknownEditVariable
	}

	delta := value - info.constant
	info.constant = value
	s.edits[variable] = info

	if r, ok := s.rows[info.tag.marker]; ok {
		if r.AddConstant(-delta) < 0 {
			s.infeasibleRows = append(s.infeasibleRows, info.tag.marker)
		}
		s.rows[info.tag.marker] = r
	} else if r, ok := s.rows[info.tag.other]; ok {
		if r.AddConstant(delta) < 0 {
			s.infeasibleRows = append(s.infeasibleRows, info.tag.other)
		}
		s.rows[info.tag.other] = r
	} else {
		for rowSymbol, current := range s.rows {
			coefficient := current.CoefficientFor(info.tag.marker)
			diff := delta * coefficient
			if diff != 0 && rowSymbol.Kind() == symbolExternal {
				s.varChanged(s.varForSymbol[rowSymbol])
			}
			if coefficient != 0 && current.AddConstant(diff) < 0 && rowSymbol.Kind() != symbolExternal {
				s.infeasibleRows = append(s.infeasibleRows, rowSymbol)
			}
			s.rows[rowSymbol] = current
		}
	}

	if err := s.dualOptimize(); err != nil {
		return ErrInternalSolver
	}
	return nil
}

// FetchChanges returns variable changes since the previous fetch.
func (s *Solver) FetchChanges() []Change {
	if s.shouldClearChanges {
		clear(s.changed)
		s.shouldClearChanges = false
	} else {
		s.shouldClearChanges = true
	}

	s.publicChanges = s.publicChanges[:0]
	for v := range s.changed {
		data, ok := s.varData[v]
		if !ok {
			continue
		}
		newValue := 0.0
		if r, ok := s.rows[data.symbol]; ok {
			newValue = r.constant
		}
		if data.value != newValue {
			s.publicChanges = append(s.publicChanges, Change{Variable: v, Value: newValue})
			data.value = newValue
			s.varData[v] = data
		}
	}

	changes := make([]Change, len(s.publicChanges))
	copy(changes, s.publicChanges)
	return changes
}

func (s *Solver) getVarSymbol(v Variable) symbol {
	data, ok := s.varData[v]
	if !ok {
		newSymbol := newSymbol(s.idTick, symbolExternal)
		s.idTick++
		s.varForSymbol[newSymbol] = v
		data = varData{value: math.NaN(), symbol: newSymbol}
	}
	data.count++
	s.varData[v] = data
	return data.symbol
}

func (s *Solver) createRow(constraint Constraint) (row, tag) {
	expr := constraint.Expression()
	r := newRow(expr.Constant())

	for _, term := range expr.Terms() {
		if nearZero(term.Coefficient()) {
			continue
		}
		termSymbol := s.getVarSymbol(term.Var())
		if other, ok := s.rows[termSymbol]; ok {
			r.InsertRow(&other, term.Coefficient())
		} else {
			r.InsertSymbol(termSymbol, term.Coefficient())
		}
	}

	constraintTag := tag{marker: invalidSymbol(), other: invalidSymbol()}
	switch constraint.Operator() {
	case LessOrEqual, GreaterOrEqual:
		coefficient := 1.0
		if constraint.Operator() == GreaterOrEqual {
			coefficient = -1.0
		}
		slack := newSymbol(s.idTick, symbolSlack)
		s.idTick++
		r.InsertSymbol(slack, coefficient)
		constraintTag.marker = slack
		if constraint.Strength().Less(Required) {
			errSymbol := newSymbol(s.idTick, symbolError)
			s.idTick++
			r.InsertSymbol(errSymbol, -coefficient)
			s.objective.InsertSymbol(errSymbol, constraint.Strength().Value())
			constraintTag.other = errSymbol
		}
	case Equal:
		if constraint.Strength().Less(Required) {
			errPlus := newSymbol(s.idTick, symbolError)
			s.idTick++
			errMinus := newSymbol(s.idTick, symbolError)
			s.idTick++
			r.InsertSymbol(errPlus, -1)
			r.InsertSymbol(errMinus, 1)
			s.objective.InsertSymbol(errPlus, constraint.Strength().Value())
			s.objective.InsertSymbol(errMinus, constraint.Strength().Value())
			constraintTag = tag{marker: errPlus, other: errMinus}
		} else {
			dummy := newSymbol(s.idTick, symbolDummy)
			s.idTick++
			r.InsertSymbol(dummy, 1)
			constraintTag.marker = dummy
		}
	}

	if r.constant < 0 {
		r.ReverseSign()
	}
	return r, constraintTag
}

func chooseSubject(r *row, constraintTag tag) symbol {
	for candidate := range r.cells {
		if candidate.Kind() == symbolExternal {
			return candidate
		}
	}
	if (constraintTag.marker.Kind() == symbolSlack || constraintTag.marker.Kind() == symbolError) &&
		r.CoefficientFor(constraintTag.marker) < 0 {
		return constraintTag.marker
	}
	if (constraintTag.other.Kind() == symbolSlack || constraintTag.other.Kind() == symbolError) &&
		r.CoefficientFor(constraintTag.other) < 0 {
		return constraintTag.other
	}
	return invalidSymbol()
}

func (s *Solver) addWithArtificialVariable(r *row) (bool, error) {
	art := newSymbol(s.idTick, symbolSlack)
	s.idTick++
	s.rows[art] = cloneRow(r)
	artificial := cloneRow(r)
	s.artificial = &artificial

	if err := s.optimize(optimizeArtificial); err != nil {
		s.artificial = nil
		return false, err
	}
	success := nearZero(s.artificial.constant)
	s.artificial = nil

	if artRow, ok := s.rows[art]; ok {
		delete(s.rows, art)
		if len(artRow.cells) == 0 {
			return success, nil
		}
		entering := anyPivotableSymbol(&artRow)
		if entering.Kind() == symbolInvalid {
			return false, nil
		}
		artRow.SolveForSymbols(art, entering)
		s.substitute(entering, &artRow)
		s.rows[entering] = artRow
	}

	for key, current := range s.rows {
		current.RemoveSymbol(art)
		s.rows[key] = current
	}
	s.objective.RemoveSymbol(art)
	return success, nil
}

func (s *Solver) substitute(substitution symbol, substituteRow *row) {
	for rowSymbol, current := range s.rows {
		constantChanged := current.Substitute(substitution, substituteRow)
		if rowSymbol.Kind() == symbolExternal && constantChanged {
			s.varChanged(s.varForSymbol[rowSymbol])
		}
		if rowSymbol.Kind() != symbolExternal && current.constant < 0 {
			s.infeasibleRows = append(s.infeasibleRows, rowSymbol)
		}
		s.rows[rowSymbol] = current
	}
	s.objective.Substitute(substitution, substituteRow)
	if s.artificial != nil {
		s.artificial.Substitute(substitution, substituteRow)
	}
}

func (s *Solver) optimize(target optimizeTarget) error {
	for {
		objective := &s.objective
		if target == optimizeArtificial {
			objective = s.artificial
		}
		entering := getEnteringSymbol(objective)
		if entering.Kind() == symbolInvalid {
			return nil
		}
		leaving, r, ok := s.getLeavingRow(entering)
		if !ok {
			return ErrInternalSolver
		}
		r.SolveForSymbols(leaving, entering)
		s.substitute(entering, &r)
		if entering.Kind() == symbolExternal && r.constant != 0 {
			s.varChanged(s.varForSymbol[entering])
		}
		s.rows[entering] = r
	}
}

func (s *Solver) dualOptimize() error {
	for len(s.infeasibleRows) > 0 {
		last := len(s.infeasibleRows) - 1
		leaving := s.infeasibleRows[last]
		s.infeasibleRows = s.infeasibleRows[:last]

		r, ok := s.rows[leaving]
		if !ok || r.constant >= 0 {
			continue
		}
		delete(s.rows, leaving)

		entering := s.getDualEnteringSymbol(&r)
		if entering.Kind() == symbolInvalid {
			return ErrInternalSolver
		}

		r.SolveForSymbols(leaving, entering)
		s.substitute(entering, &r)
		if entering.Kind() == symbolExternal && r.constant != 0 {
			s.varChanged(s.varForSymbol[entering])
		}
		s.rows[entering] = r
	}
	return nil
}

func getEnteringSymbol(objective *row) symbol {
	for candidate, value := range objective.cells {
		if candidate.Kind() != symbolDummy && value < 0 {
			return candidate
		}
	}
	return invalidSymbol()
}

func (s *Solver) getDualEnteringSymbol(r *row) symbol {
	ratio := math.Inf(1)
	entering := invalidSymbol()
	for candidate, value := range r.cells {
		if value > 0 && candidate.Kind() != symbolDummy {
			tempRatio := s.objective.CoefficientFor(candidate) / value
			if tempRatio < ratio {
				ratio = tempRatio
				entering = candidate
			}
		}
	}
	return entering
}

func anyPivotableSymbol(r *row) symbol {
	for candidate := range r.cells {
		if candidate.Kind() == symbolSlack || candidate.Kind() == symbolError {
			return candidate
		}
	}
	return invalidSymbol()
}

func (s *Solver) getLeavingRow(entering symbol) (symbol, row, bool) {
	ratio := math.Inf(1)
	found := invalidSymbol()
	for rowSymbol, r := range s.rows {
		if rowSymbol.Kind() == symbolExternal {
			continue
		}
		temp := r.CoefficientFor(entering)
		if temp < 0 {
			tempRatio := -r.constant / temp
			if tempRatio < ratio {
				ratio = tempRatio
				found = rowSymbol
			}
		}
	}
	if found.Kind() == symbolInvalid {
		return invalidSymbol(), row{}, false
	}
	r := s.rows[found]
	delete(s.rows, found)
	return found, r, true
}

func (s *Solver) getMarkerLeavingRow(marker symbol) (symbol, row, bool) {
	r1 := math.Inf(1)
	r2 := math.Inf(1)
	first := invalidSymbol()
	second := invalidSymbol()
	third := invalidSymbol()

	for rowSymbol, r := range s.rows {
		coefficient := r.CoefficientFor(marker)
		if coefficient == 0 {
			continue
		}
		if rowSymbol.Kind() == symbolExternal {
			third = rowSymbol
		} else if coefficient < 0 {
			ratio := -r.constant / coefficient
			if ratio < r1 {
				r1 = ratio
				first = rowSymbol
			}
		} else {
			ratio := r.constant / coefficient
			if ratio < r2 {
				r2 = ratio
				second = rowSymbol
			}
		}
	}

	found := first
	if found.Kind() == symbolInvalid {
		found = second
	}
	if found.Kind() == symbolInvalid {
		found = third
	}
	if found.Kind() == symbolInvalid {
		return invalidSymbol(), row{}, false
	}
	if found.Kind() == symbolExternal && s.rows[found].constant != 0 {
		s.varChanged(s.varForSymbol[found])
	}
	r := s.rows[found]
	delete(s.rows, found)
	return found, r, true
}

func (s *Solver) removeConstraintEffects(constraint Constraint, constraintTag tag) {
	if constraintTag.marker.Kind() == symbolError {
		s.removeMarkerEffects(constraintTag.marker, constraint.Strength().Value())
	}
	if constraintTag.other.Kind() == symbolError {
		s.removeMarkerEffects(constraintTag.other, constraint.Strength().Value())
	}
}

func (s *Solver) removeMarkerEffects(marker symbol, strength float64) {
	if r, ok := s.rows[marker]; ok {
		s.objective.InsertRow(&r, -strength)
	} else {
		s.objective.InsertSymbol(marker, -strength)
	}
}

func allDummies(r *row) bool {
	for candidate := range r.cells {
		if candidate.Kind() != symbolDummy {
			return false
		}
	}
	return true
}

func cloneRow(r *row) row {
	clone := newRow(r.constant)
	maps.Copy(clone.cells, r.cells)
	return clone
}

func (s *Solver) snapshotState() solverStateSnapshot {
	snapshot := solverStateSnapshot{
		constraints:        make(map[Constraint]tag, len(s.constraints)),
		edits:              make(map[Variable]editInfo, len(s.edits)),
		varData:            make(map[Variable]varData, len(s.varData)),
		varForSymbol:       make(map[symbol]Variable, len(s.varForSymbol)),
		publicChanges:      append([]Change(nil), s.publicChanges...),
		changed:            make(map[Variable]struct{}, len(s.changed)),
		shouldClearChanges: s.shouldClearChanges,
		rows:               make(map[symbol]row, len(s.rows)),
		infeasibleRows:     append([]symbol(nil), s.infeasibleRows...),
		objective:          cloneRow(&s.objective),
		idTick:             s.idTick,
	}
	maps.Copy(snapshot.constraints, s.constraints)
	maps.Copy(snapshot.edits, s.edits)
	maps.Copy(snapshot.varData, s.varData)
	maps.Copy(snapshot.varForSymbol, s.varForSymbol)
	maps.Copy(snapshot.changed, s.changed)
	for rowSymbol, current := range s.rows {
		rowClone := cloneRow(&current)
		snapshot.rows[rowSymbol] = rowClone
	}
	if s.artificial != nil {
		artificial := cloneRow(s.artificial)
		snapshot.artificial = &artificial
	}
	return snapshot
}

func (s *Solver) restoreState(snapshot solverStateSnapshot) {
	s.constraints = snapshot.constraints
	s.edits = snapshot.edits
	s.varData = snapshot.varData
	s.varForSymbol = snapshot.varForSymbol
	s.publicChanges = snapshot.publicChanges
	s.changed = snapshot.changed
	s.shouldClearChanges = snapshot.shouldClearChanges
	s.rows = snapshot.rows
	s.infeasibleRows = snapshot.infeasibleRows
	s.objective = snapshot.objective
	s.artificial = snapshot.artificial
	s.idTick = snapshot.idTick
}

func (s *Solver) varChanged(v Variable) {
	if s.shouldClearChanges {
		clear(s.changed)
		s.shouldClearChanges = false
	}
	s.changed[v] = struct{}{}
}
