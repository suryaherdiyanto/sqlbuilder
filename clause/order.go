package clause

import (
	"fmt"
	"strings"
)

type OrderField struct {
	Field     string
	Direction OrderDirection
}

type Order struct {
	OrderingFields []OrderField
}

func (o Order) Parse(dialect SQLDialector) string {
	if len(o.OrderingFields) == 0 {
		return ""
	}

	stmt := "ORDER BY "
	for i, orderField := range o.OrderingFields {
		stmt += fmt.Sprintf("%s %s", orderField.Field, strings.ToUpper(string(orderField.Direction)))
		if i < len(o.OrderingFields)-1 {
			stmt += ", "
		}
	}

	return stmt
}
