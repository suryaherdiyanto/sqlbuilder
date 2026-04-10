package clause

import (
	"fmt"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

type WhereDay struct {
	Field string
	Op    Operator
	Conj  Conjuction
	Value any
}

func (w WhereDay) Parse(d SQLDialector) string {
	dialectName := d.GetName()
	col := fmt.Sprintf("%s%s%s", d.GetColumnQuoteLeft(), w.Field, d.GetColumnQuoteRight())

	switch dialectName {
	case dialect.MySQL:
		return fmt.Sprintf("DAY(%s) %s ?", col, w.Op)
	case dialect.PostgreSQL:
		return fmt.Sprintf("EXTRACT(DAY FROM %s) %s ?", col, w.Op)
	case dialect.SQLite:
		return fmt.Sprintf("strftime('%%d', %s) %s ?", col, w.Op)
	default:
		return ""
	}
}
