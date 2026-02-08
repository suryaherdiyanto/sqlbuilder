package clause

type Update struct {
	Table string
	Rows  map[string]any
	WhereStatements
}

func (u Update) Parse(dialect SQLDialector) string {
	return dialect.ParseUpdate(u)
}
