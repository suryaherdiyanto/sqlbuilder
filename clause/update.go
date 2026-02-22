package clause

import (
	"fmt"
	"slices"
	"strings"
)

type Update struct {
	Table  string
	Rows   map[string]any
	Values []any
}

func (u Update) Parse(dialect SQLDialector) (string, Update) {
	stmt := fmt.Sprintf("UPDATE %s SET ", u.Table)
	keys := make([]string, 0, len(u.Rows))

	for k := range u.Rows {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, k := range keys {
		stmt += fmt.Sprintf("%s%s%s = ?, ", dialect.GetColumnQuoteLeft(), k, dialect.GetColumnQuoteRight())
		if val, ok := u.Rows[k]; ok {
			u.Values = append(u.Values, val)
		}
	}

	stmt = strings.TrimRight(stmt, ", ")

	return stmt, u
}

func (u Update) GetArguments() []any {
	return u.Values
}
