package clause

import "github.com/suryaherdiyanto/sqlbuilder/pkg"

type GroupBy struct {
	Fields []string
}

func (g GroupBy) Parse(dialect SQLDialector) string {
	if len(g.Fields) == 0 {
		return ""
	}

	stmt := " GROUP BY "
	for i, field := range g.Fields {
		stmt += pkg.ColumnSplitter(field, dialect.GetColumnQuoteLeft(), dialect.GetColumnQuoteRight())
		if i < len(g.Fields)-1 {
			stmt += ","
		}
	}

	return stmt
}
