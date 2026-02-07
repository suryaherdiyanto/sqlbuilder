package clause

type OrderField struct {
	Field     string
	Direction OrderDirection
}

type Order struct {
	OrderingFields []OrderField
}

func (o Order) Parse(dialect SQLDialector) string {
	return dialect.ParseOrder(o)
}
