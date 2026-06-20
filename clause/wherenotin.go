package clause

import "fmt"

type WhereNotIn struct {
	Field        string
	Values       []any
	Conj         Conjuction
	SubStatement Select
}

func (wi WhereNotIn) Parse(d SQLDialector) string {
	if wi.SubStatement.Table != "" {
		subStmt, _ := wi.SubStatement.Parse(d)
		return fmt.Sprintf("%s%s%s NOT IN (%s)", d.GetColumnQuoteLeft(), wi.Field, d.GetColumnQuoteRight(), subStmt)

	}
	inValues := ""

	for i := range wi.Values {
		inValues += d.GetDelimiter()

		if i < len(wi.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s%s%s NOT IN(%s)", d.GetColumnQuoteLeft(), wi.Field, d.GetColumnQuoteRight(), inValues)
}
