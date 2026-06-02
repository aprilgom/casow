package casow

type symbolKind int

const (
	symbolInvalid symbolKind = iota
	symbolExternal
	symbolSlack
	symbolError
	symbolDummy
)

type symbol struct {
	id   uint64
	kind symbolKind
}

func newSymbol(id uint64, kind symbolKind) symbol {
	return symbol{id: id, kind: kind}
}

func invalidSymbol() symbol {
	return symbol{id: 0, kind: symbolInvalid}
}

func (s symbol) Kind() symbolKind {
	return s.kind
}
