package clause

type Select struct {
	Table   string
	Joins   []Join
	Columns []string
	Values  []any
}

func (s Select) Parse(dialect SQLDialector) (string, Select) {
	return dialect.ParseSelect(s)
}

func (s Select) GetArguments() []any {
	return s.Values
}
