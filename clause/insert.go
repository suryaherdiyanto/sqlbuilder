package clause

type Insert struct {
	Table string
	Rows  []map[string]any
}

func (in Insert) Parse(dialect SQLDialector) string {
	return dialect.ParseInsert(in)
}
