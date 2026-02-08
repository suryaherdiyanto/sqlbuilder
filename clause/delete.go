package clause

type Delete struct {
	Table string
	WhereStatements
}

func (d Delete) Parse(dialect SQLDialector) string {
	return dialect.ParseDelete(d)
}
