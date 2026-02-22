package clause

import "fmt"

type WhereNotBetween struct {
	Field string
	Start any
	End   any
	Conj  Conjuction
}

func (w WhereNotBetween) Parse(dialect SQLDialector) string {
	return fmt.Sprintf("%s%s%s NOT BETWEEN %s AND %s", dialect.GetColumnQuoteLeft(), w.Field, dialect.GetColumnQuoteRight(), dialect.GetDelimiter(), dialect.GetDelimiter())
}
