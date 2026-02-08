package clause

type WhereNotBetween struct {
	Field string
	Start any
	End   any
	Conj  Conjuction
}

func (w WhereNotBetween) Parse(dialect SQLDialector) string {
	return dialect.ParseWhereNotBetween(w)
}
