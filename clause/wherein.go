package clause

type WhereIn struct {
	Field        string
	Values       []any
	Conj         Conjuction
	SubStatement Select
}

func (w WhereIn) Parse(dialect SQLDialector) string {
	return dialect.ParseWhereIn(w)
}
