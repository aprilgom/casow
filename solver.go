package casow

type Solver struct{}

type Change struct {
	Variable Variable
	Value    float64
}

func NewSolver() *Solver {
	return &Solver{}
}

func (s *Solver) AddConstraint(constraint Constraint) error {
	return ErrInternalSolver
}

func (s *Solver) RemoveConstraint(constraint Constraint) error {
	return ErrInternalSolver
}

func (s *Solver) FetchChanges() []Change {
	return nil
}
