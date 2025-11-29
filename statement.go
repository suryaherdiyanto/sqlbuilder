package sqlbuilder

import (
	"fmt"
	"strings"
)

type Operator string
type JoinType string
type OrderDirection string

const (
	OperatorEqual             Operator = "="
	OperatorLessThan                   = "<"
	OperatorLessThanEqual              = "<="
	OperatorGreaterThan                = ">"
	OperatorGreatherThanEqual          = ">="
	OperatorNot                        = "!="
)

const (
	LeftJoin  JoinType = "left join"
	RightJoin          = "right join"
	InnerJoin          = "inner join"
)

const (
	OrderDirectionASC  OrderDirection = "asc"
	OrderDirectionDESC                = "desc"
)

type WhereParser interface {
	Parse() string
}

type WhereInParser interface {
	Parse() string
}

type JoinParser interface {
	Parse() string
}

type StatementParser interface {
	Parse() string
}

type OrderParser interface {
	Parse() string
}

type Where struct {
	Field string
	Op    Operator
	Value any
}

func (w *Where) Parse() string {
	val := w.Value
	switch w.Value.(type) {
	case string:
		val = fmt.Sprintf("'%s'", w.Value)
	}

	return fmt.Sprintf("%s %s %v", w.Field, w.Op, val)
}

type WhereIn struct {
	Field  string
	Values []any
}

func (w *WhereIn) Parse() string {
	inValues := ""
	for i, v := range w.Values {
		switch v.(type) {
		case string:
			inValues += fmt.Sprintf("'%s'", v)
		default:
			inValues += fmt.Sprintf("%v", v)
		}

		if i < len(w.Values)-1 {
			inValues += ","
		}
	}
	return fmt.Sprintf("%s IN(%s)", w.Field, inValues)
}

type Join struct {
	Type       JoinType
	OtherTable string
	On         JoinON
}

type JoinON struct {
	LeftTable  string
	LeftValue  any
	RightValue any
	RightTable string
}

type Order struct {
	Field     string
	Direction OrderDirection
}

func (o *Order) Parse() string {
	return fmt.Sprintf("ORDER BY %s %s", o.Field, o.Direction)
}

type SelectStatement struct {
	Table           string
	Columns         []string
	WhereStatements []Where
	WhereIn         WhereIn
	Joins           []Join
	SubQueries      []SubQuery
	Ordering        Order
	setStatement    string
}

func (s *SelectStatement) Parse() string {
	stmt := `SELECT %s FROM %s WHERE `

	fields := strings.Join(s.Columns, ",")

	for _, v := range s.WhereStatements {
		stmt += v.Parse()
	}

	if s.Ordering != (Order{}) {
		stmt += s.Ordering.Parse()
	}

	return fmt.Sprintf(stmt, fields, s.Table)
}

type SubQuery struct {
	SubStatement SelectStatement
}
