package clause

import "fmt"

type Offset struct {
	Count int64
}

func (o Offset) Parse(dialect SQLDialector) string {
	if o.Count == 0 {
		return ""
	}

	return fmt.Sprintf(" OFFSET %s", dialect.GetDelimiter())
}
