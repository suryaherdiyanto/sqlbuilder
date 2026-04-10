package clause

import (
	"fmt"
	"strings"

	"github.com/suryaherdiyanto/sqlbuilder/pkg"
)

type JoinON struct {
	Operator   Operator
	LeftField  string
	RightField string
}

type Join struct {
	Type        JoinType
	SecondTable string
	On          JoinON
}

func (j Join) Parse(dialect SQLDialector) string {
	leftField := pkg.ColumnSplitter(j.On.LeftField, dialect.GetColumnQuoteLeft(), dialect.GetColumnQuoteRight())
	rightField := pkg.ColumnSplitter(j.On.RightField, dialect.GetColumnQuoteLeft(), dialect.GetColumnQuoteRight())
	rightTable := dialect.GetColumnQuoteLeft() + j.SecondTable + dialect.GetColumnQuoteRight()

	return fmt.Sprintf("%s %s ON %s %s %s", strings.ToUpper(string(j.Type)), rightTable, leftField, j.On.Operator, rightField)
}
