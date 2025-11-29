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

type SelectStatement struct {
	Table               string
	Columns             []string
	WhereStatements     []Where
	WhereInStatement    WhereIn
	WhereNotInStatement WhereNotIn
	Joins               []Join
	SubQueries          []SubQuery
	Ordering            Order
	Limit               int64
	Offset              int64
	setStatement        string
}

type WhereIn struct {
	Field  string
	Values []any
}

type WhereNotIn struct {
	Field  string
	Values []any
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

func (w *Where) Parse() string {
	val := w.Value
	switch w.Value.(type) {
	case string:
		val = fmt.Sprintf("'%s'", w.Value)
	}

	return fmt.Sprintf("%s %s %v", w.Field, w.Op, val)
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

func (w *WhereNotIn) Parse() string {
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
	return fmt.Sprintf("%s NOT IN(%s)", w.Field, inValues)

}

func (o *Order) Parse() string {
	return fmt.Sprintf("ORDER BY %s %s", o.Field, o.Direction)
}

func (s *SelectStatement) Parse() string {
	stmt := `SELECT %s FROM %s WHERE `

	fields := strings.Join(s.Columns, ",")

	for _, v := range s.WhereStatements {
		stmt += v.Parse()
	}

	if s.WhereInStatement.Field != "" || len(s.WhereInStatement.Values) > 0 {
		if len(s.WhereStatements) > 0 {
			stmt += " AND "
		}
		stmt += s.WhereInStatement.Parse()
	}

	if s.WhereNotInStatement.Field != "" || len(s.WhereNotInStatement.Values) > 0 {
		if len(s.WhereStatements) > 0 || s.WhereInStatement.Field != "" {
			stmt += " AND "
		}
		stmt += s.WhereNotInStatement.Parse()
	}

	if s.Ordering != (Order{}) {
		stmt += s.Ordering.Parse()
	}

	if s.Limit > 0 {
		stmt += fmt.Sprintf(" LIMIT %d", s.Limit)
	}

	if s.Offset > 0 {
		stmt += fmt.Sprintf(" OFFSET %d", s.Offset)
	}

	return fmt.Sprintf(stmt, fields, s.Table)
}

type SubQuery struct {
	SubStatement SelectStatement
}
