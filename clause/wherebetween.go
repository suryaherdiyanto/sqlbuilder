package clause

import "fmt"

type WhereBetween struct {
	Field string
	Start any
	End   any
	Conj  Conjuction
}

func (w WhereBetween) Parse(dialect SQLDialector) string {
	return fmt.Sprintf("%s%s%s BETWEEN %s AND %s", dialect.GetColumnQuoteLeft(), w.Field, dialect.GetColumnQuoteRight(), dialect.GetDelimiter(), dialect.GetDelimiter())
}
