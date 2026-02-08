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
	if w.SubStatement.Table != "" {
		subStmt, _ := w.SubStatement.Parse(d)
		return fmt.Sprintf("%s%s%s %s (%s)", d.ColumnQuoteLeft, w.Field, d.ColumnQuoteRight, w.Op, subStmt)
	}

	if strings.Contains(w.Field, ".") {
		return fmt.Sprintf("%s %s %s", columnSplitter(w.Field, d.ColumnQuoteLeft, d.ColumnQuoteRight), w.Op, d.Delimiter)
	}

	return fmt.Sprintf("%s%s%s %s %s", d.ColumnQuoteLeft, w.Field, d.ColumnQuoteRight, w.Op, d.Delimiter)
}

func (d *SQLiteDialect) ParseWhereBetween(wb clause.WhereBetween) string {
	return fmt.Sprintf("%s%s%s BETWEEN %s AND %s", d.ColumnQuoteLeft, wb.Field, d.ColumnQuoteRight, d.Delimiter, d.Delimiter)
}

func (d *SQLiteDialect) ParseWhereNotBetween(wb clause.WhereNotBetween) string {
	return fmt.Sprintf("%s%s%s NOT BETWEEN %s AND %s", d.ColumnQuoteLeft, wb.Field, d.ColumnQuoteRight, d.Delimiter, d.Delimiter)
}

func (d *SQLiteDialect) ParseWhereIn(wi clause.WhereIn) string {
	if wi.SubStatement.Table != "" {
		subStmt, _ := wi.SubStatement.Parse(d)
		return fmt.Sprintf("%s%s%s IN (%s)", d.ColumnQuoteLeft, wi.Field, d.ColumnQuoteRight, subStmt)

	}
	inValues := ""

	for i := range wi.Values {
		inValues += d.Delimiter

		if i < len(wi.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s%s%s IN(%s)", d.ColumnQuoteLeft, wi.Field, d.ColumnQuoteRight, inValues)
}

func (d *SQLiteDialect) ParseWhereNotIn(wi clause.WhereNotIn) string {
	if wi.SubStatement.Table != "" {
		subStmt, _ := wi.SubStatement.Parse(d)
		return fmt.Sprintf("%s%s%s NOT IN (%s)", d.ColumnQuoteLeft, wi.Field, d.ColumnQuoteRight, subStmt)

	}
	inValues := ""

	for i := range wi.Values {
		inValues += d.Delimiter

		if i < len(wi.Values)-1 {
			inValues += ","
		}
	}

	return fmt.Sprintf("%s%s%s NOT IN(%s)", d.ColumnQuoteLeft, wi.Field, d.ColumnQuoteRight, inValues)
}

func (d *SQLiteDialect) ParseJoin(j clause.Join) string {
	return fmt.Sprintf("%s %s ON %s.%s %s %s.%s", strings.ToUpper(string(j.Type)), d.ColumnQuoteLeft+j.SecondTable+d.ColumnQuoteRight, d.ColumnQuoteLeft+j.FirstTable+d.ColumnQuoteRight, d.ColumnQuoteLeft+j.On.LeftField+d.ColumnQuoteRight, j.On.Operator, d.ColumnQuoteLeft+j.SecondTable+d.ColumnQuoteRight, d.ColumnQuoteLeft+j.On.RightField+d.ColumnQuoteRight)
}

func (d *SQLiteDialect) ParseGroup(g clause.GroupBy) string {
	if len(g.Fields) == 0 {
		return ""
	}

	stmt := " GROUP BY "
	for i, field := range g.Fields {
		stmt += columnSplitter(field, d.ColumnQuoteLeft, d.ColumnQuoteRight)
		if i < len(g.Fields)-1 {
			stmt += ", "
		}
	}

	return stmt
}

func (d *SQLiteDialect) ParseOrder(o clause.Order) string {
	if len(o.OrderingFields) == 0 {
		return ""
	}

	stmt := "ORDER BY "
	for i, orderField := range o.OrderingFields {
		stmt += fmt.Sprintf("%s %s", orderField.Field, strings.ToUpper(string(orderField.Direction)))
		if i < len(o.OrderingFields)-1 {
			stmt += ", "
		}
	}

	return stmt
}

func (d *SQLiteDialect) ParseJoins(j []clause.Join) string {
	stmt := ""
	for _, v := range j {
		stmt += fmt.Sprintf(" %s", d.ParseJoin(v))
	}

	return stmt
}

func (d *SQLiteDialect) ParseSelect(s clause.Select) (string, clause.Select) {
	stmt := `SELECT %s FROM %s%s%s`

	columns := []string{}
	for _, col := range s.Columns {
		columns = append(columns, columnSplitter(col, d.ColumnQuoteLeft, d.ColumnQuoteRight))
	}

	fields := strings.Join(columns, ",")

	stmt += d.ParseJoins(s.Joins)

	if len(s.WhereStatements.Where) > 0 || len(s.WhereStatements.WhereIn) > 0 || len(s.WhereStatements.WhereNotIn) > 0 || len(s.WhereStatements.WhereBetween) > 0 || len(s.WhereStatements.WhereNotBetween) > 0 {
		stmt += " WHERE "
	}

	stmt += d.ParseWhereStatements(&s.WhereStatements)

	stmt += d.ParseWhereInStatements(&s.WhereStatements)

	stmt += d.ParseWhereNotInStatements(&s.WhereStatements)

	stmt += d.ParseWhereBetweenStatements(&s.WhereStatements)

	stmt += d.ParseWhereNotBetweenStatements(&s.WhereStatements)

	stmt += d.ParseGroup(s.GroupBy)

	stmt += d.ParseOrder(s.Order)

	stmt += d.ParseLimit(s.Limit)

	stmt += d.ParseOffset(s.Offset)

	return fmt.Sprintf(stmt, fields, d.ColumnQuoteLeft, s.Table, d.ColumnQuoteRight), s
}

func (d *SQLiteDialect) ParseLimit(l clause.Limit) string {
	if l.Count == 0 {
		return ""
	}

	return fmt.Sprintf(" LIMIT %s", d.Delimiter)
}

func (d *SQLiteDialect) ParseOffset(o clause.Offset) string {
	if o.Count == 0 {
		return ""
	}

	return fmt.Sprintf(" OFFSET %s", d.Delimiter)
}

func (d *SQLiteDialect) ParseWhereStatements(ws *clause.WhereStatements) string {
	stmt := ""
	for i, v := range ws.Where {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse(d)
		if v.Value != nil {
			ws.Values = append(ws.Values, v.Value)
		}
	}
	return stmt
}

func (d *SQLiteDialect) ParseWhereInStatements(ws *clause.WhereStatements) string {
	stmt := ""
	if len(ws.Where) > 0 {
		for _, v := range ws.WhereIn {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
	}

	if len(ws.WhereIn) > 0 {
		for _, v := range ws.WhereIn {
			stmt += v.Parse(d)
			ws.Values = append(ws.Values, v.Values...)
		}
	}
	return stmt
}

func (d *SQLiteDialect) ParseWhereNotInStatements(ws *clause.WhereStatements) string {
	stmt := ""
	for i, v := range ws.WhereNotIn {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse(d)
		ws.Values = append(ws.Values, v.Values...)
	}
	return stmt
}

func (d *SQLiteDialect) ParseWhereBetweenStatements(ws *clause.WhereStatements) string {
	stmt := ""
	for i, v := range ws.WhereBetween {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse(d)
		ws.Values = append(ws.Values, v.Start, v.End)
	}
	return stmt
}

func (d *SQLiteDialect) ParseWhereNotBetweenStatements(ws *clause.WhereStatements) string {
	stmt := ""
	for i, v := range ws.WhereNotBetween {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse(d)
		ws.Values = append(ws.Values, v.Start, v.End)
	}
	return stmt
}

func (d *SQLiteDialect) ParseUpdate(u clause.Update) string {
	stmt := fmt.Sprintf("UPDATE %s SET ", u.Table)
	keys := make([]string, 0, len(u.Rows))

	for k := range u.Rows {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, k := range keys {
		stmt += fmt.Sprintf("`%s` = ?, ", k)
		if val, ok := u.Rows[k]; ok {
			u.Values = append(u.Values, val)
		}
	}

	stmt = strings.TrimRight(stmt, ", ")

	stmt += " WHERE "
	stmt += d.ParseWhereStatements(&u.WhereStatements)
	stmt += d.ParseWhereInStatements(&u.WhereStatements)
	stmt += d.ParseWhereNotInStatements(&u.WhereStatements)
	stmt += d.ParseWhereBetweenStatements(&u.WhereStatements)
	stmt += d.ParseWhereNotBetweenStatements(&u.WhereStatements)
	return stmt
}

func (s *SQLiteDialect) ParseDelete(d clause.Delete) string {
	stmt := fmt.Sprintf("DELETE FROM %s", s.ColumnQuoteLeft+d.Table+s.ColumnQuoteRight)

	stmt += " WHERE "
	stmt += s.ParseWhereStatements(&d.WhereStatements)
	stmt += s.ParseWhereInStatements(&d.WhereStatements)
	stmt += s.ParseWhereNotInStatements(&d.WhereStatements)
	stmt += s.ParseWhereBetweenStatements(&d.WhereStatements)
	stmt += s.ParseWhereNotBetweenStatements(&d.WhereStatements)

	return stmt
}

func NewSQLiteDialect() *SQLiteDialect {
	return &SQLiteDialect{Delimiter: "?", ColumnQuoteLeft: "`", ColumnQuoteRight: "`"}
}
