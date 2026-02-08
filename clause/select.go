package clause

type Select struct {
	Table   string
	Columns []string
	Joins   []Join
	GroupBy GroupBy
	Order   Order
	Limit   Limit
	Offset  Offset
	WhereStatements
	Values []any
}

func (s Select) Parse(dialect SQLDialector) (string, Select) {
	return dialect.ParseSelect(s)
}

func (s Select) GetArguments() []any {
	s.Values = append(s.Values, s.WhereStatements.Values...)
	return s.Values
}
