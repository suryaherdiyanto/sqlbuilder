package clause

import "fmt"

type WhereIn struct {
	Field        string
	Values       []any
	Conj         Conjuction
	SubStatement SubStatement
}

func (wi WhereIn) Parse(d SQLDialector) string {
	if wi.SubStatement.Table != "" {
		subStmt, _ := wi.SubStatement.Select.Parse(d)
		return fmt.Sprintf("%s%s%s IN (%s)", d.GetColumnQuoteLeft(), wi.Field, d.GetColumnQuoteRight(), subStmt)

	}
	inValues := ""

	for i := range wi.Values {
		inValues += d.GetDelimiter()

		if i < len(wi.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s%s%s IN(%s)", d.GetColumnQuoteLeft(), wi.Field, d.GetColumnQuoteRight(), inValues)
}
