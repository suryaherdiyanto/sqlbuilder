package clause

import "fmt"

type Limit struct {
	Count int64
}

func (l Limit) Parse(dialect SQLDialector) string {
	if l.Count == 0 {
		return ""
	}

	return fmt.Sprintf(" LIMIT %s", dialect.GetDelimiter())
}
