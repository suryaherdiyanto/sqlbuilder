package clause

import (
	"fmt"
	"slices"
	"strings"
)

type Insert struct {
	Table string
	Rows  []map[string]any
}

func (in Insert) Parse(dialect SQLDialector) string {
	columns := ""
	keys := []string{}

	if len(in.Rows) > 0 {
		keys = make([]string, 0, len(in.Rows[0]))

		for k := range in.Rows[0] {
			keys = append(keys, k)
		}
		slices.Sort(keys)

		for _, k := range keys {
			columns += fmt.Sprintf("%s%s%s,", dialect.GetColumnQuoteLeft(), k, dialect.GetColumnQuoteRight())
		}
	}

	columns = strings.TrimRight(columns, ",")

	insertValues := ""
	for i := range len(in.Rows) {
		rowValues := ""
		for idx := range keys {
			rowValues += dialect.GetDelimiter()
			if idx < len(keys)-1 {
				rowValues += ","
			}
		}

		insertValues += fmt.Sprintf("(%s)", rowValues)
		if i < len(in.Rows)-1 {
			insertValues += ","
		}
	}

	return fmt.Sprintf("INSERT INTO %s(%s) VALUES%s", in.Table, columns, insertValues)
}
