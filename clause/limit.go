package clause

type Limit struct {
	Count int
}

func (l Limit) Parse(dialect SQLDialector) string {
	return dialect.ParseLimit(l)
}
