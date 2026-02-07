package clause

type Offset struct {
	Count int
}

func (o Offset) Parse(dialect SQLDialector) string {
	return dialect.ParseOffset(o)
}
