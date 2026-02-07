package clause

type WhereBetween struct {
	Field string
	Start any
	End   any
	Conj  Conjuction
}

func (w WhereBetween) Parse(dialect SQLDialector) string {
	return dialect.ParseWhereBetween(w)
}
