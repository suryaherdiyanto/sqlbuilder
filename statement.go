package sqlbuilder

import (
	"fmt"
	"strings"
)

type SelectStatement struct {
	Table              string
	Columns            []string
	JoinStatements     []Join
	GroupByStatement   GroupBy
	Ordering           Order
	Limit              int64
	Offset             int64
	HasExistsClause    bool
	HasNotExistsClause bool
	WhereStatements
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

type DeleteStatement struct {
	Table                     string
	WhereStatements           []Where
	WhereInStatements         []WhereIn
	WhereNotInStatements      []WhereNotIn
	WhereBetweenStatements    []WhereBetween
	WhereNotBetweenStatements []WhereNotBetween
	Values                    []any
}

type WhereStatements struct {
	Where           []Where
	WhereIn         []WhereIn
	WhereNotIn      []WhereNotIn
	WhereBetween    []WhereBetween
	WhereNotBetween []WhereNotBetween
	Values          []any
}

type InsertStatement struct {
	Table string
	Rows  []map[string]any
}

func (ws *WhereStatements) ParseWheres() string {
	stmt := ""
	for i, v := range ws.Where {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse()
		if v.Value != nil {
			ws.Values = append(ws.Values, v.Value)
		}
	}
	return stmt
}

func (ws *WhereStatements) ParseWhereIn() string {
	stmt := ""
	if len(ws.Where) > 0 {
		for _, v := range ws.WhereIn {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
	}

	if len(ws.WhereIn) > 0 {
		for _, v := range ws.WhereIn {
			stmt += v.Parse()
			ws.Values = append(ws.Values, v.Values...)
		}
	}

	return stmt
}

func (ws *WhereStatements) ParseWhereNotIn() string {
	stmt := ""
	if len(ws.Where) > 0 {
		for _, v := range ws.WhereNotIn {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
	}

	if len(ws.WhereNotIn) > 0 {
		for _, v := range ws.WhereNotIn {
			stmt += v.Parse()
			ws.Values = append(ws.Values, v.Values...)
		}
	}

	return stmt
}

func (ws *WhereStatements) ParseWhereBetweens() string {
	stmt := ""
	for i, v := range ws.WhereBetween {
		if i >= 1 || len(ws.WhereBetween) > 0 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
		stmt += v.Parse()
		ws.Values = append(ws.Values, v.Start, v.End)
	}

	return stmt
}

func (ws *WhereStatements) ParseWhereNotBetweens() string {
	stmt := ""
	for i, v := range ws.WhereNotBetween {
		if i >= 1 || len(ws.WhereNotBetween) > 0 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
		stmt += v.Parse()
		ws.Values = append(ws.Values, v.Start, v.End)
	}

	return stmt
}

func (ws *WhereStatements) ParseAllWheres() string {
	stmt := ""
	if len(ws.Where) > 0 || len(ws.WhereBetween) > 0 || len(ws.WhereNotBetween) > 0 || len(ws.WhereIn) > 0 || len(ws.WhereNotIn) > 0 {
		stmt += " WHERE "
	}

	stmt += ws.ParseWheres()
	stmt += ws.ParseWhereIn()
	stmt += ws.ParseWhereNotIn()
	stmt += ws.ParseWhereBetweens()
	stmt += ws.ParseWhereNotBetweens()
	return stmt
}

func (s *SelectStatement) ParseJoins() string {
	stmt := ""
	for _, v := range s.JoinStatements {
		stmt += fmt.Sprintf(" %s", v.Parse())
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

func (s *SelectStatement) Parse() string {
	stmt := `SELECT %s FROM %s`

	fields := strings.Join(s.Columns, ",")

	stmt += s.ParseJoins()

	stmt += s.WhereStatements.ParseAllWheres()

	stmt += s.ParseGroupings()

	stmt += s.ParseOrdering()

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

func (s *SelectStatement) GetArguments() []any {
	return s.Values
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

func (s *UpdateStatement) GetArguments() []any {
	return s.Values
}

func (d *DeleteStatement) Parse() string {
	stmt := fmt.Sprintf("DELETE FROM %s", d.Table)

	if len(d.WhereStatements) > 0 || len(d.WhereBetweenStatements) > 0 || len(d.WhereNotBetweenStatements) > 0 || len(d.WhereInStatements) > 0 || len(d.WhereNotInStatements) > 0 {
		stmt += " WHERE "
	}

	stmt += d.ParseWheres()

	stmt += d.ParseWhereBetweens()

	stmt += d.ParseWhereNotBetweens()

	stmt += d.ParseWhereIn()

	stmt += d.ParseWhereNotIn()

	return stmt
}

func (d *DeleteStatement) ParseWheres() string {
	stmt := ""
	for i, v := range d.WhereStatements {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse()
		d.Values = append(d.Values, v.Value)
	}
	return stmt
}

func (d *DeleteStatement) ParseWhereIn() string {
	stmt := ""
	if len(d.WhereStatements) > 0 {
		for _, v := range d.WhereInStatements {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
	}

	if len(d.WhereInStatements) > 0 {
		for _, v := range d.WhereInStatements {
			stmt += v.Parse()
			d.Values = append(d.Values, v.Values...)
		}
	}

	return stmt
}

func (d *DeleteStatement) ParseWhereNotIn() string {
	stmt := ""
	if len(d.WhereStatements) > 0 {
		for _, v := range d.WhereNotInStatements {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
	}

	if len(d.WhereNotInStatements) > 0 {
		for _, v := range d.WhereNotInStatements {
			stmt += v.Parse()
			d.Values = append(d.Values, v.Values...)
		}
	}

	return stmt
}

func (d *DeleteStatement) ParseWhereBetweens() string {
	stmt := ""
	for i, v := range d.WhereBetweenStatements {
		if i >= 1 || len(d.WhereStatements) > 0 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
		stmt += v.Parse()
		d.Values = append(d.Values, v.Start, v.End)
	}

	return stmt
}

func (d *DeleteStatement) ParseWhereNotBetweens() string {
	stmt := ""
	for i, v := range d.WhereNotBetweenStatements {
		if i >= 1 || len(d.WhereStatements) > 0 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
		stmt += v.Parse()
		d.Values = append(d.Values, v.Start, v.End)
	}

	return stmt
}

func (d *DeleteStatement) GetArguments() []any {
	return d.Values
}
