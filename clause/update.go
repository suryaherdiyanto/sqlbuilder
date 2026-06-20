package clause

import (
	"fmt"
	"slices"
	"strings"

	"github.com/suryaherdiyanto/sqlbuilder/dialect"
)

type Update struct {
	Table  string
	Rows   map[string]any
	Values []any
}

func (u Update) Parse(d SQLDialector, i int) (string, Update) {
	stmt := fmt.Sprintf("UPDATE %s SET ", u.Table)
	keys := make([]string, 0, len(u.Rows))

	for k := range u.Rows {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for j, k := range keys {
		delimiter := d.GetDelimiter()
		if d.GetName() == dialect.PostgreSQL {
			delimiter = fmt.Sprintf("$%d", i+j)
		}

		stmt += fmt.Sprintf("%s%s%s = %s, ", d.GetColumnQuoteLeft(), k, d.GetColumnQuoteRight(), delimiter)
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
