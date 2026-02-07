package dialect

import (
	"fmt"
	"slices"
	"strings"

	"github.com/suryaherdiyanto/sqlbuilder/clause"
)

type SQLiteDialect struct {
	Delimiter        string
	ColumnQuoteLeft  string
	ColumnQuoteRight string
}

func (d *SQLiteDialect) ParseInsert(in clause.Insert) string {
	columns := ""
	values := ""

	if len(in.Rows) > 0 {
		keys := make([]string, 0, len(in.Rows[0]))

		for k := range in.Rows[0] {
			keys = append(keys, k)
		}
		slices.Sort(keys)

		for _, k := range keys {
			columns += fmt.Sprintf("%s%s%s,", d.ColumnQuoteLeft, k, d.ColumnQuoteRight)
			values += "?,"
		}
	}

	columns = strings.TrimRight(columns, ",")
	values = strings.TrimRight(values, ",")

	insertValues := ""
	for i := range len(in.Rows) {
		insertValues += fmt.Sprintf("(%s)", values)
		if i < len(in.Rows)-1 {
			insertValues += ","
		}
	}

	return fmt.Sprintf("INSERT INTO %s(%s) VALUES%s", in.Table, columns, insertValues)
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
	return fmt.Sprintf("%s %s ON %s.%s %s %s.%s", strings.ToUpper(string(j.Type)), d.ColumnQuoteLeft+j.SecondTable+d.ColumnQuoteRight, d.ColumnQuoteLeft+j.FirstTable+d.ColumnQuoteRight, d.ColumnQuoteLeft+j.On.LeftField+d.ColumnQuoteRight, j.On.Operator, d.ColumnQuoteLeft+j.SecondTable+d.ColumnQuoteRight, d.ColumnQuoteLeft+j.On.RightField+d.ColumnQuoteRight)
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
	return fmt.Sprintf("LIMIT %s", d.Delimiter)
}

func (d *SQLiteDialect) ParseOffset(o clause.Offset) string {
	return fmt.Sprintf("OFFSET %s", d.Delimiter)
}

func NewSQLiteDialect() *SQLiteDialect {
	return &SQLiteDialect{Delimiter: "?", ColumnQuoteLeft: "`", ColumnQuoteRight: "`"}
}
