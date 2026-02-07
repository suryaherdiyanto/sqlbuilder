package clause

type WhereBetween struct {
	Field string
	Start any
	End   any
	Conj  Conjuction
}
