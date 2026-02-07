package clause

type WhereNotBetween struct {
	Field string
	Start any
	End   any
	Conj  Conjuction
}
