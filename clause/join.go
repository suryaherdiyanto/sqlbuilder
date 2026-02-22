package clause

import (
	"fmt"
	"strings"
)

type JoinON struct {
	Operator   Operator
	LeftField  string
	RightField string
}

type Join struct {
	Type        JoinType
	FirstTable  string
	SecondTable string
	On          JoinON
}

func (j Join) Parse(dialect SQLDialector) string {
	return fmt.Sprintf("%s %s ON %s.%s %s %s.%s", strings.ToUpper(string(j.Type)), dialect.GetColumnQuoteLeft()+j.SecondTable+dialect.GetColumnQuoteRight(), dialect.GetColumnQuoteLeft()+j.FirstTable+dialect.GetColumnQuoteRight(), dialect.GetColumnQuoteLeft()+j.On.LeftField+dialect.GetColumnQuoteRight(), j.On.Operator, dialect.GetColumnQuoteLeft()+j.SecondTable+dialect.GetColumnQuoteRight(), dialect.GetColumnQuoteLeft()+j.On.RightField+dialect.GetColumnQuoteRight())
}
