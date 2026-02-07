package dialect

import (
	"fmt"
	"strings"

	"github.com/suryaherdiyanto/sqlbuilder/clause"
)

type SQLiteDialect struct {
	Where           clause.Where
	WhereIn         clause.WhereIn
	WhereNotIn      clause.WhereNotIn
	WhereBetween    clause.WhereBetween
	WhereNotBetween clause.WhereNotBetween
	GroupBy         clause.GroupBy
	OrderBy         clause.Order
	Join            clause.Join
}

func (d *SQLiteDialect) ParseWhere() string {
	return fmt.Sprintf("%s %s ?", d.Where.Field, d.Where.Op)
}

func (d *SQLiteDialect) ParseWhereBetween() string {
	return fmt.Sprintf("%s BETWEEN ? AND ?", d.WhereBetween.Field)
}

func (d *SQLiteDialect) ParseWhereNotBetween() string {
	return fmt.Sprintf("%s NOT BETWEEN ? AND ?", d.WhereNotBetween.Field)
}

func (d *SQLiteDialect) ParseWhereIn() string {
	inValues := ""

	for i := range d.WhereIn.Values {
		inValues += "?"

		if i < len(d.WhereIn.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s IN(%s)", d.WhereIn.Field, inValues)
}

func (d *SQLiteDialect) ParseWhereNotIn() string {
	inValues := ""

	for i := range d.WhereNotIn.Values {
		inValues += "?"

		if i < len(d.WhereNotIn.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s NOT IN(%s)", d.WhereNotIn.Field, inValues)
}

func (d *SQLiteDialect) ParseJoin() string {
	return fmt.Sprintf("%s %s ON %s.%v %s %s.%v", strings.ToUpper(string(d.Join.Type)), d.Join.SecondTable, d.Join.FirstTable, d.Join.On.LeftValue, d.Join.On.Operator, d.Join.SecondTable, d.Join.On.RightValue)
}

func (d *SQLiteDialect) ParseGroup() string {
	stmt := "GROUP BY "
	for i, field := range d.GroupBy.Fields {
		stmt += field
		if i < len(d.GroupBy.Fields)-1 {
			stmt += ", "
		}
	}

	return stmt
}

func (d *SQLiteDialect) ParseOrder() string {
	stmt := "ORDER BY "
	for i, orderField := range d.OrderBy.OrderingFields {
		stmt += fmt.Sprintf("%s %s", orderField.Field, strings.ToUpper(string(orderField.Direction)))
		if i < len(d.OrderBy.OrderingFields)-1 {
			stmt += ", "
		}
	}

	return stmt
}

func (d *SQLiteDialect) NewWhere(field string, op clause.Operator, value any) {
	d.Where = clause.Where{
		Field: field,
		Op:    op,
		Value: value,
	}
}

func (d *SQLiteDialect) NewWhereIn(field string, values []any, conj clause.Conjuction) {
	d.WhereIn = clause.WhereIn{
		Field:  field,
		Values: values,
		Conj:   conj,
	}
}

func (d *SQLiteDialect) NewWhereNotIn(field string, values []any, conj clause.Conjuction) {
	d.WhereNotIn = clause.WhereNotIn{
		Field:  field,
		Values: values,
		Conj:   conj,
	}
}

func (d *SQLiteDialect) NewWhereBetween(field string, start, end any, conj clause.Conjuction) {
	d.WhereBetween = clause.WhereBetween{
		Field: field,
		Start: start,
		End:   end,
		Conj:  conj,
	}
}

func (d *SQLiteDialect) NewWhereNotBetween(field string, start, end any, conj clause.Conjuction) {
	d.WhereNotBetween = clause.WhereNotBetween{
		Field: field,
		Start: start,
		End:   end,
		Conj:  conj,
	}
}

func (d *SQLiteDialect) NewJoin(joinType clause.JoinType, firstTable, secondTable string, on clause.JoinON) {
	d.Join = clause.Join{
		Type:        joinType,
		FirstTable:  firstTable,
		SecondTable: secondTable,
		On:          on,
	}
}

func (d *SQLiteDialect) NewOrder(orderFields []clause.OrderField) {
	d.OrderBy = clause.Order{
		OrderingFields: orderFields,
	}
}

func (d *SQLiteDialect) NewGroup(fields []string) {
	d.GroupBy = clause.GroupBy{
		Fields: fields,
	}
}

func NewSQLiteDialect() *SQLiteDialect {
	return &SQLiteDialect{}
}
