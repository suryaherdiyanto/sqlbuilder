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

type CrossJoin struct {
	SecondTable string
}

func (j Join) Parse(d SQLDialector) string {
	leftField := pkg.ColumnSplitter(j.On.LeftField, d.GetColumnQuoteLeft(), d.GetColumnQuoteRight())
	rightField := pkg.ColumnSplitter(j.On.RightField, d.GetColumnQuoteLeft(), d.GetColumnQuoteRight())
	rightTable := d.GetColumnQuoteLeft() + j.SecondTable + d.GetColumnQuoteRight()

	return fmt.Sprintf("%s %s ON %s %s %s", strings.ToUpper(string(j.Type)), rightTable, leftField, j.On.Operator, rightField)
}

func (j CrossJoin) Parse(d SQLDialector) string {
	rightTable := d.GetColumnQuoteLeft() + j.SecondTable + d.GetColumnQuoteRight()

	return fmt.Sprintf("%s %s", strings.ToUpper(string(CrossJoinType)), rightTable)
}
