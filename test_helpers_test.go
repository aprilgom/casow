package casow

type testValues struct {
	values map[Variable]float64
}

func newValues() (func(Variable) float64, func([]Change)) {
	values := &testValues{values: make(map[Variable]float64)}

	valueOf := func(variable Variable) float64 {
		return values.values[variable]
	}
	updateValues := func(changes []Change) {
		for _, change := range changes {
			values.values[change.Variable] = change.Value
		}
	}

	return valueOf, updateValues
}
