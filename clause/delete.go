package clause

type Delete struct {
	Table  string
	Values []any
}

func (d Delete) Parse(dialect SQLDialector) (string, Delete) {
	return dialect.ParseDelete(d)
}

func (d Delete) GetArguments() []any {
	return d.Values
}
