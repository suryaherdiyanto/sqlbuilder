package clause

import (
	"fmt"
)

type WhereNotIn struct {
	Field  string
	Values []any
	Conj   Conjuction
}

func (w *WhereNotIn) Parse() string {
	inValues := ""

	for i, _ := range w.Values {
		inValues += "?"

		if i < len(w.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s NOT IN(%s)", w.Field, inValues)

}
