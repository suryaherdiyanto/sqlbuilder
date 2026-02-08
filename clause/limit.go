package clause

type Limit struct {
	Count int64
}

func (l Limit) Parse(dialect SQLDialector) string {
	return dialect.ParseLimit(l)
}
