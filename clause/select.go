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
}

func (s Select) Parse(dialect SQLDialector) string {
	return dialect.ParseSelect(s)
}
