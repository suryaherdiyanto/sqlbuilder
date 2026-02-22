package clause

import "fmt"

type WhereIn struct {
	Field        string
	Values       []any
	Conj         Conjuction
	SubStatement SubStatement
}

func (wi WhereIn) Parse(dialect SQLDialector) string {
	if wi.SubStatement.Table != "" {
		subStmt, _ := wi.SubStatement.Select.Parse(dialect)
		subWhereStmt := wi.SubStatement.WhereStatements.Parse(dialect)
		return fmt.Sprintf("%s%s%s IN (%s%s)", dialect.GetColumnQuoteLeft(), wi.Field, dialect.GetColumnQuoteRight(), subStmt, subWhereStmt)

	}
	inValues := ""

	for i := range wi.Values {
		inValues += dialect.GetDelimiter()

		if i < len(wi.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s%s%s IN(%s)", dialect.GetColumnQuoteLeft(), wi.Field, dialect.GetColumnQuoteRight(), inValues)
}
