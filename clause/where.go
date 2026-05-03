package clause

import (
	"fmt"
	"strings"

	"github.com/suryaherdiyanto/sqlbuilder/pkg"
)

type SubStatement struct {
	Select
}

type Where struct {
	Field        string
	Op           Operator
	Value        any
	Conj         Conjuction
	SubStatement SubStatement
}

func (w Where) Parse(dialect SQLDialector) string {
	field := w.GetField(dialect)
	if strings.Contains(w.Field, ".") {
		field = pkg.ColumnSplitter(w.Field, dialect.GetColumnQuoteLeft(), dialect.GetColumnQuoteRight())
	}

	if w.SubStatement.Table != "" {
		subStmt, _ := w.SubStatement.Select.Parse(dialect)
		if w.Op == OperatorExists {
			return fmt.Sprintf("%s (%s)", w.Op, subStmt)
		}
		return fmt.Sprintf("%s %s (%s)", field, w.Op, subStmt)
	}

	return fmt.Sprintf("%s %s %s", field, w.Op, dialect.GetDelimiter())
}

func (w Where) GetField(dialect SQLDialector) string {
	return fmt.Sprintf("%s%s%s", dialect.GetColumnQuoteLeft(), w.Field, dialect.GetColumnQuoteRight())
}
