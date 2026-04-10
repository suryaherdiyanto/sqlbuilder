package clause

import (
	"fmt"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

type WhereDate struct {
	Field string
	Op    Operator
	Conj  Conjuction
	Value any
}

func (w WhereDate) Parse(d SQLDialector) string {
	dialectName := d.GetName()
	col := fmt.Sprintf("%s%s%s", d.GetColumnQuoteLeft(), w.Field, d.GetColumnQuoteRight())

	switch dialectName {
	case dialect.MySQL:
		return fmt.Sprintf("DATE(%s) %s ?", col, w.Op)
	case dialect.PostgreSQL:
		return fmt.Sprintf("CAST(%s AS DATE) %s ?", col, w.Op)
	case dialect.SQLite:
		return fmt.Sprintf("DATE(%s) %s ?", col, w.Op)
	default:
		return ""
	}
}
