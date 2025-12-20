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
	OperatorExists                     = "EXISTS"
	OperatorNotExists                  = "NOT EXISTS"
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
	GroupByStatement          GroupBy
	Ordering                  Order
	Limit                     int64
	Offset                    int64
	Values                    []any
	HasExistsClause           bool
	HasNotExistsClause        bool
}

type WhereIn struct {
	Field        string
	Values       []any
	Conj         Conjuction
	SubStatement SelectStatement
}

type WhereNotIn struct {
	Field        string
	Values       []any
	Conj         Conjuction
	SubStatement SelectStatement
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
	Type        JoinType
	FirstTable  string
	SecondTable string
	On          JoinON
}

type JoinON struct {
	Operator   Operator
	LeftValue  any
	RightValue any
}

type OrderField struct {
	Field     string
	Direction OrderDirection
}

type Order struct {
	OrderingFields []OrderField
}

type GroupBy struct {
	Fields []string
}

type InsertStatement struct {
	Table string
	Rows  []map[string]any
}

type UpdateStatement struct {
	Table                     string
	Rows                      map[string]any
	WhereStatements           []Where
	WhereInStatements         []WhereIn
	WhereNotInStatements      []WhereNotIn
	WhereBetweenStatements    []WhereBetween
	WhereNotBetweenStatements []WhereNotBetween
	Values                    []any
}

func (s *SelectStatement) ParseWheres() string {
	stmt := ""
	for i, v := range s.WhereStatements {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse()
		s.Values = append(s.Values, v.Value)
	}
	return stmt
}

func (w *Where) Parse() string {

	if (!reflect.DeepEqual(w.SubStatement, SelectStatement{})) {
		val := fmt.Sprintf("(%s)", w.SubStatement.Parse())
		if w.Op == OperatorExists || w.Op == OperatorNotExists {
			return fmt.Sprintf("%s %v", w.Op, val)
		}

		return fmt.Sprintf("%s %s %s", w.Field, w.Op, val)
	}

	return fmt.Sprintf("%s %s ?", w.Field, w.Op)
}

func (w *WhereIn) Parse() string {
	inValues := ""
	if !reflect.DeepEqual(w.SubStatement, SelectStatement{}) {
		return fmt.Sprintf("%s IN (%s)", w.Field, w.SubStatement.Parse())
	}

	for i, _ := range w.Values {
		inValues += "?"

		if i < len(w.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s IN(%s)", w.Field, inValues)
}

func (w *WhereNotIn) Parse() string {
	inValues := ""
	if !reflect.DeepEqual(w.SubStatement, SelectStatement{}) {
		return fmt.Sprintf("%s NOT IN (%s)", w.Field, w.SubStatement.Parse())
	}

	for i, _ := range w.Values {
		inValues += "?"

		if i < len(w.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s NOT IN(%s)", w.Field, inValues)

}

func (wb *WhereBetween) Parse() string {
	return fmt.Sprintf("%s BETWEEN ? AND ?", wb.Field)
}

func (wnb *WhereNotBetween) Parse() string {
	return fmt.Sprintf("%s NOT BETWEEN ? AND ?", wnb.Field)
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

func (g *GroupBy) Parse() string {
	stmt := ""
	for i, field := range g.Fields {
		stmt += field
		if i < len(g.Fields)-1 {
			stmt += ","
		}
	}

	return stmt
}

func (s *SelectStatement) ParseJoins() string {
	stmt := ""
	for _, v := range s.JoinStatements {
		stmt += fmt.Sprintf(" %s", v.Parse())
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
		s.Values = append(s.Values, v.Start, v.End)
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
		s.Values = append(s.Values, v.Start, v.End)
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
			s.Values = append(s.Values, v.Values...)
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
			s.Values = append(s.Values, v.Values...)
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

func (s *SelectStatement) ParseGroupings() string {
	stmt := ""
	if len(s.GroupByStatement.Fields) > 0 {
		stmt = fmt.Sprintf(" GROUP BY %s", s.GroupByStatement.Parse())
	}

	return stmt
}

func (j *Join) Parse() string {
	return fmt.Sprintf("%s %s ON %s.%v %s %s.%v", strings.ToUpper(string(j.Type)), j.SecondTable, j.FirstTable, j.On.LeftValue, j.On.Operator, j.SecondTable, j.On.RightValue)
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

	stmt += s.ParseGroupings()

	if s.Limit > 0 {
		stmt += " LIMIT ?"
		s.Values = append(s.Values, s.Limit)
	}

	if s.Offset > 0 {
		stmt += " OFFSET ?"
		s.Values = append(s.Values, s.Offset)
	}

	return fmt.Sprintf(stmt, fields, s.Table)
}

func (si *InsertStatement) Parse() string {
	columns := ""
	values := ""

	if len(si.Rows) > 0 {
		for k, _ := range si.Rows[0] {
			columns += "`" + k + "`" + ","
			values += "?,"
		}
	}
	columns = strings.TrimRight(columns, ",")
	values = strings.TrimRight(values, ",")

	stmt := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", si.Table, columns, values)

	return stmt
}

func (su *UpdateStatement) Parse() string {
	stmt := fmt.Sprintf("UPDATE %s SET ", su.Table)
	for k, v := range su.Rows {
		stmt += fmt.Sprintf("`%s` = ?, ", k)
		su.Values = append(su.Values, v)
	}

	stmt = strings.TrimRight(stmt, ", ")

	if len(su.WhereStatements) > 0 || len(su.WhereBetweenStatements) > 0 || len(su.WhereNotBetweenStatements) > 0 || len(su.WhereInStatements) > 0 || len(su.WhereNotInStatements) > 0 {
		stmt += " WHERE "
	}

	stmt += su.ParseWheres()

	stmt += su.ParseWhereBetweens()

	stmt += su.ParseWhereNotBetweens()

	stmt += su.ParseWhereIn()

	stmt += su.ParseWhereNotIn()

	return stmt
}

func (s *UpdateStatement) ParseWheres() string {
	stmt := ""
	for i, v := range s.WhereStatements {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse()
		s.Values = append(s.Values, v.Value)
	}
	return stmt
}

func (s *UpdateStatement) ParseWhereIn() string {
	stmt := ""
	if len(s.WhereStatements) > 0 {
		for _, v := range s.WhereInStatements {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
	}

	if len(s.WhereInStatements) > 0 {
		for _, v := range s.WhereInStatements {
			stmt += v.Parse()
			s.Values = append(s.Values, v.Values...)
		}
	}

	return stmt
}

func (s *UpdateStatement) ParseWhereNotIn() string {
	stmt := ""
	if len(s.WhereStatements) > 0 {
		for _, v := range s.WhereNotInStatements {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
	}

	if len(s.WhereNotInStatements) > 0 {
		for _, v := range s.WhereNotInStatements {
			stmt += v.Parse()
			s.Values = append(s.Values, v.Values...)
		}
	}

	return stmt
}

func (s *UpdateStatement) ParseWhereBetweens() string {
	stmt := ""
	for i, v := range s.WhereBetweenStatements {
		if i >= 1 || len(s.WhereStatements) > 0 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
		stmt += v.Parse()
		s.Values = append(s.Values, v.Start, v.End)
	}

	return stmt
}

func (s *UpdateStatement) ParseWhereNotBetweens() string {
	stmt := ""
	for i, v := range s.WhereNotBetweenStatements {
		if i >= 1 || len(s.WhereStatements) > 0 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
		stmt += v.Parse()
		s.Values = append(s.Values, v.Start, v.End)
	}

	return stmt
}
