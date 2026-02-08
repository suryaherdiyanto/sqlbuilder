package clause

type Offset struct {
	Count int64
}

func (o Offset) Parse(dialect SQLDialector) string {
	return dialect.ParseOffset(o)
}
