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

type WhereIn struct {
	Field        string
	Values       []any
	Conj         Conjuction
	SubStatement SelectStatement
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

type WhereNotIn struct {
	Field        string
	Values       []any
	Conj         Conjuction
	SubStatement SelectStatement
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

type WhereBetween struct {
	Field string
	Start any
	End   any
	Conj  Conjuction
}

func (wb *WhereBetween) Parse() string {
	return fmt.Sprintf("%s BETWEEN ? AND ?", wb.Field)
}

type WhereNotBetween struct {
	Field string
	Start any
	End   any
	Conj  Conjuction
}

func (wnb *WhereNotBetween) Parse() string {
	return fmt.Sprintf("%s NOT BETWEEN ? AND ?", wnb.Field)
}

type Join struct {
	Type        JoinType
	FirstTable  string
	SecondTable string
	On          JoinON
}

func (j *Join) Parse() string {
	return fmt.Sprintf("%s %s ON %s.%v %s %s.%v", strings.ToUpper(string(j.Type)), j.SecondTable, j.FirstTable, j.On.LeftValue, j.On.Operator, j.SecondTable, j.On.RightValue)
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

type GroupBy struct {
	Fields []string
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
