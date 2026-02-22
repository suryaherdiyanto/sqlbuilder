package clause

import "fmt"

type WhereNotIn struct {
	Field        string
	Values       []any
	Conj         Conjuction
	SubStatement Select
}

func (wi WhereNotIn) Parse(dialect SQLDialector) string {
	if wi.SubStatement.Table != "" {
		subStmt, _ := wi.SubStatement.Parse(dialect)
		return fmt.Sprintf("%s%s%s NOT IN (%s)", dialect.GetColumnQuoteLeft(), wi.Field, dialect.GetColumnQuoteRight(), subStmt)

	}
	inValues := ""

	for i := range wi.Values {
		inValues += dialect.GetDelimiter()

		if i < len(wi.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s%s%s NOT IN(%s)", dialect.GetColumnQuoteLeft(), wi.Field, dialect.GetColumnQuoteRight(), inValues)
}
