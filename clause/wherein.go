package clause

type WhereIn struct {
	Field  string
	Values []any
	Conj   Conjuction
}

func (w WhereIn) Parse(dialect SQLDialector) string {
	return dialect.ParseWhereIn(w)
}
