package clause

import (
	"fmt"
	"strings"

	"github.com/suryaherdiyanto/sqlbuilder/pkg"
)

type Select struct {
	Table   string
	Joins   []Join
	Columns []string
	Values  []any
}

func (s Select) Parse(dialect SQLDialector) (string, Select) {
	stmt := `SELECT %s FROM %s`

	columns := []string{}
	for _, col := range s.Columns {
		if col == "*" {
			columns = append(columns, "*")
			continue
		}

		columns = append(columns, pkg.ColumnSplitter(col, dialect.GetColumnQuoteLeft(), dialect.GetColumnQuoteRight()))
	}

	fields := strings.Join(columns, ",")

	return fmt.Sprintf(stmt, fields, s.Table), s
}

func (s Select) ParseJoins(dialect SQLDialector) string {
	stmt := ""
	for _, v := range s.Joins {
		stmt += fmt.Sprintf(" %s", v.Parse(dialect))
	}

	return stmt
}

func (s Select) GetArguments() []any {
	return s.Values
}
