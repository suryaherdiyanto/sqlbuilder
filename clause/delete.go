package clause

import "fmt"

type Delete struct {
	Table  string
	Values []any
}

func (d Delete) Parse(dialect SQLDialector) (string, Delete) {
	stmt := fmt.Sprintf("DELETE FROM %s", d.Table)

	return stmt, d
}

func (d Delete) GetArguments() []any {
	return d.Values
}
