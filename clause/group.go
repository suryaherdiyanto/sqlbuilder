package clause

type GroupBy struct {
	Fields []string
}

func (g GroupBy) Parse(dialect SQLDialector) string {
	return dialect.ParseGroup(g)
}
