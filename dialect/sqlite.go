package dialect

import (
	"fmt"
	"strings"

	"github.com/suryaherdiyanto/sqlbuilder/clause"
)

type SQLiteDialect struct {
	Delimiter        string
	ColumnQuoteLeft  string
	ColumnQuoteRight string
}

func (d *SQLiteDialect) ParseWhere(w clause.Where) string {
	return fmt.Sprintf("%s %s %s", w.Field, w.Op, d.Delimiter)
}

func (d *SQLiteDialect) ParseWhereBetween(wb clause.WhereBetween) string {
	return fmt.Sprintf("%s BETWEEN %s AND %s", wb.Field, d.Delimiter, d.Delimiter)
}

func (d *SQLiteDialect) ParseWhereNotBetween(wb clause.WhereNotBetween) string {
	return fmt.Sprintf("%s NOT BETWEEN %s AND %s", wb.Field, d.Delimiter, d.Delimiter)
}

func (d *SQLiteDialect) ParseWhereIn(wi clause.WhereIn) string {
	inValues := ""

	for i := range wi.Values {
		inValues += d.Delimiter

		if i < len(wi.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s IN(%s)", wi.Field, inValues)
}

func (d *SQLiteDialect) ParseWhereNotIn(wi clause.WhereNotIn) string {
	inValues := ""

	for i := range wi.Values {
		inValues += d.Delimiter

		if i < len(wi.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s NOT IN(%s)", wi.Field, inValues)
}

func (d *SQLiteDialect) ParseJoin(j clause.Join) string {
	return fmt.Sprintf("%s %s ON %s.%v %s %s.%v", strings.ToUpper(string(j.Type)), j.SecondTable, j.FirstTable, j.On.LeftValue, j.On.Operator, j.SecondTable, j.On.RightValue)
}

func (d *SQLiteDialect) ParseGroup(g clause.GroupBy) string {
	stmt := "GROUP BY "
	for i, field := range g.Fields {
		stmt += field
		if i < len(g.Fields)-1 {
			stmt += ", "
		}
	}

	return stmt
}

func (d *SQLiteDialect) ParseOrder(o clause.Order) string {
	stmt := "ORDER BY "
	for i, orderField := range o.OrderingFields {
		stmt += fmt.Sprintf("%s %s", orderField.Field, strings.ToUpper(string(orderField.Direction)))
		if i < len(o.OrderingFields)-1 {
			stmt += ", "
		}
	}

	return stmt
}

func (d *SQLiteDialect) ParseLimit(l clause.Limit) string {
	return fmt.Sprintf("LIMIT ?", d.Delimiter)
}

func (d *SQLiteDialect) ParseOffset(o clause.Offset) string {
	return fmt.Sprintf("OFFSET ?", d.Delimiter)
}

func NewSQLiteDialect() *SQLiteDialect {
	return &SQLiteDialect{Delimiter: "?", ColumnQuoteLeft: "`", ColumnQuoteRight: "`"}
}
