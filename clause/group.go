package clause

import "github.com/suryaherdiyanto/sqlbuilder/pkg"

type GroupBy struct {
	Fields []string
}

func (g GroupBy) Parse(d SQLDialector) string {
	if len(g.Fields) == 0 {
		return ""
	}

	stmt := "GROUP BY "
	for i, field := range g.Fields {
		stmt += pkg.ColumnSplitter(field, d.GetColumnQuoteLeft(), d.GetColumnQuoteRight())
		if i < len(g.Fields)-1 {
			stmt += ","
		}
	}

	return stmt
}
