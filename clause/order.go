package clause

import "fmt"

type OrderField struct {
	Field     string
	Direction OrderDirection
}

type Order struct {
	OrderingFields []OrderField
}

func (o *Order) Parse() string {
	stmt := "ORDER BY "
	for i, f := range o.OrderingFields {
		stmt += fmt.Sprintf("%s %s", f.Field, f.Direction)
		if i < len(o.OrderingFields)-1 {
			stmt += ", "
		}
	}

	return stmt
}
