package clause

type WhereGroup struct {
	Conj   Conjuction
	Wheres []Where
}

type Where struct {
	Field  string
	Op     Operator
	Value  any
	Conj   Conjuction
	Groups []WhereGroup
}

type WhereStatements struct {
	Where           []Where
	WhereIn         []WhereIn
	WhereNotIn      []WhereNotIn
	WhereBetween    []WhereBetween
	WhereNotBetween []WhereNotBetween
	Values          []any
}

func (w Where) Parse(dialect SQLDialector) string {
	return dialect.ParseWhere(w)
}

func (w WhereStatements) Parse(dialect SQLDialector) string {
	return dialect.ParseWhereStatements(w)
}
