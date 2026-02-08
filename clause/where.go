package clause

type WhereGroup struct {
	Conj   Conjuction
	Wheres []Where
}

type SubStatement struct {
	Select
	WhereStatements
}

type Where struct {
	Field        string
	Op           Operator
	Value        any
	Conj         Conjuction
	Groups       []WhereGroup
	SubStatement SubStatement
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

func (w *WhereStatements) Parse(dialect SQLDialector) string {
	stmt := ""
	if len(w.Where) > 0 || len(w.WhereIn) > 0 || len(w.WhereNotIn) > 0 || len(w.WhereBetween) > 0 || len(w.WhereNotBetween) > 0 {
		stmt += " WHERE "
	}

	stmt += dialect.ParseWhereStatements(w)

	stmt += dialect.ParseWhereInStatements(w)

	stmt += dialect.ParseWhereNotInStatements(w)

	stmt += dialect.ParseWhereBetweenStatements(w)

	stmt += dialect.ParseWhereNotBetweenStatements(w)

	return stmt
}
