package sqlbuilder

import (
	"fmt"
	"reflect"
	"strings"
)

type Operator string
type JoinType string
type OrderDirection string
type Conjuction string

const (
	OperatorEqual             Operator = "="
	OperatorLessThan                   = "<"
	OperatorLessThanEqual              = "<="
	OperatorGreaterThan                = ">"
	OperatorGreatherThanEqual          = ">="
	OperatorNot                        = "!="
	OperatorLike                       = "LIKE"
	OperatorNotLike                    = "NOT LIKE"
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

const (
	ConjuctionAnd Conjuction = "AND"
	ConjuctionOr             = "OR"
)

type WhereGroup struct {
	Conj   Conjuction
	Wheres []Where
}

type Where struct {
	Field        string
	Op           Operator
	Value        any
	Conj         Conjuction
	Groups       []WhereGroup
	SubStatement SelectStatement
}

type SelectStatement struct {
	Table                     string
	Columns                   []string
	WhereStatements           []Where
	WhereBetweenStatements    []WhereBetween
	WhereNotBetweenStatements []WhereNotBetween
	JoinStatements            []Join
	WhereInStatements         []WhereIn
	WhereNotInStatements      []WhereNotIn
	Ordering                  Order
	Limit                     int64
	Offset                    int64
	setStatement              string
}

type WhereIn struct {
	Field  string
	Values []any
	Conj   Conjuction
}

type WhereNotIn struct {
	Field  string
	Values []any
	Conj   Conjuction
}

type WhereBetween struct {
	Field string
	Start any
	End   any
	Conj  Conjuction
}

type WhereNotBetween struct {
	Field string
	Start any
	End   any
	Conj  Conjuction
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

type OrderField struct {
	Field     string
	Direction OrderDirection
}

type Order struct {
	OrderingFields []OrderField
}

func (s *SelectStatement) ParseWheres() string {
	stmt := ""
	for i, v := range s.WhereStatements {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse()
	}
	return stmt
}

func (w *Where) Parse() string {
	val := w.Value
	switch w.Value.(type) {
	case string:
		val = fmt.Sprintf("'%s'", w.Value)
	}

	if (!reflect.DeepEqual(w.SubStatement, SelectStatement{})) {
		val = fmt.Sprintf("(%s)", w.SubStatement.Parse())
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

func (wb *WhereBetween) Parse() string {
	return fmt.Sprintf("%s BETWEEN %v AND %v", wb.Field, wb.Start, wb.End)
}

func (wnb *WhereNotBetween) Parse() string {
	return fmt.Sprintf("%s NOT BETWEEN %v AND %v", wnb.Field, wnb.Start, wnb.End)
}

func (o *Order) Parse() string {
	stmt := "ORDER BY "
	for i, f := range o.OrderingFields {
		stmt += fmt.Sprintf("%s %s", f.Field, f.Direction)
		if i < len(o.OrderingFields)-1 {
			stmt += ", "
		}
	}

	return stmt
}

func (s *SelectStatement) ParseJoins() string {
	stmt := ""
	for _, v := range s.JoinStatements {
		stmt += fmt.Sprintf("%s ", v.Parse())
	}

	return stmt
}

func (s *SelectStatement) ParseWhereBetweens() string {
	stmt := ""
	for i, v := range s.WhereBetweenStatements {
		if i >= 1 || len(s.WhereStatements) > 0 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
		stmt += v.Parse()
	}

	return stmt
}

func (s *SelectStatement) ParseWhereNotBetweens() string {
	stmt := ""
	for i, v := range s.WhereNotBetweenStatements {
		if i >= 1 || len(s.WhereStatements) > 0 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
		stmt += v.Parse()
	}

	return stmt
}

func (s *SelectStatement) ParseWhereIn() string {
	stmt := ""
	if len(s.WhereStatements) > 0 {
		for _, v := range s.WhereInStatements {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
	}

	if len(s.WhereInStatements) > 0 {
		for _, v := range s.WhereInStatements {
			stmt += v.Parse()
		}
	}

	return stmt
}

func (s *SelectStatement) ParseWhereNotIn() string {
	stmt := ""
	if len(s.WhereStatements) > 0 {
		for _, v := range s.WhereNotInStatements {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
	}

	if len(s.WhereNotInStatements) > 0 {
		for _, v := range s.WhereNotInStatements {
			stmt += v.Parse()
		}
	}

	return stmt
}

func (s *SelectStatement) ParseOrdering() string {
	stmt := ""
	if len(s.Ordering.OrderingFields) > 0 {
		stmt += s.Ordering.Parse()
	}

	return stmt
}

func (j *Join) Parse() string {
	return fmt.Sprintf("%s %s ON %s.%v = %s.%v", strings.ToUpper(string(j.Type)), j.OtherTable, j.On.LeftTable, j.On.LeftValue, j.On.RightTable, j.On.RightValue)
}

func (s *SelectStatement) Parse() string {
	stmt := `SELECT %s FROM %s`

	fields := strings.Join(s.Columns, ",")

	stmt += s.ParseJoins()

	if len(s.WhereStatements) > 0 || len(s.WhereBetweenStatements) > 0 || len(s.WhereNotBetweenStatements) > 0 || len(s.WhereInStatements) > 0 || len(s.WhereNotInStatements) > 0 {
		stmt += " WHERE "
	}

	stmt += s.ParseWheres()

	stmt += s.ParseWhereBetweens()

	stmt += s.ParseWhereNotBetweens()

	stmt += s.ParseWhereIn()

	stmt += s.ParseWhereNotIn()

	stmt += s.ParseOrdering()

	if s.Limit > 0 {
		stmt += fmt.Sprintf(" LIMIT %d", s.Limit)
	}

	if s.Offset > 0 {
		stmt += fmt.Sprintf(" OFFSET %d", s.Offset)
	}

	return fmt.Sprintf(stmt, fields, s.Table)
}
