package clause

import "fmt"

type WhereBetween struct {
	Field string
	Start any
	End   any
	Conj  Conjuction
}

func (w WhereBetween) Parse(d SQLDialector) string {
	return fmt.Sprintf("%s%s%s BETWEEN %s AND %s", d.GetColumnQuoteLeft(), w.Field, d.GetColumnQuoteRight(), d.GetDelimiter(), d.GetDelimiter())
}
