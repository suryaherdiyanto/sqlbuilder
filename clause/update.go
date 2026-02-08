package clause

type Update struct {
	Table  string
	Rows   map[string]any
	Values []any
}

func (u Update) Parse(dialect SQLDialector) (string, Update) {
	return dialect.ParseUpdate(u)
}

func (u Update) GetArguments() []any {
	return u.Values
}
