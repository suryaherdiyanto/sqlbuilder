package clause

type WhereIn struct {
	Field  string
	Values []any
	Conj   Conjuction
}
