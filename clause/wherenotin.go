package clause

type WhereNotIn struct {
	Field        string
	Values       []any
	Conj         Conjuction
	SubStatement Select
}

func (w WhereNotIn) Parse(dialect SQLDialector) string {
	return dialect.ParseWhereNotIn(w)
}
