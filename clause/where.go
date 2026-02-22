package clause

import (
	"fmt"
	"strings"

	"github.com/suryaherdiyanto/sqlbuilder/pkg"
)

type WhereGroup struct {
	Conj   Conjuction
	Wheres []Where
}

type SubStatement struct {
	Select
	WhereStatements
}

type Where struct {
	Field        string
	Op           Operator
	Value        any
	Conj         Conjuction
	Groups       []WhereGroup
	SubStatement SubStatement
}

type WhereStatements struct {
	Where           []Where
	WhereIn         []WhereIn
	WhereNotIn      []WhereNotIn
	WhereBetween    []WhereBetween
	WhereNotBetween []WhereNotBetween
	Values          []any
}

func (w Where) Parse(dialect SQLDialector) string {
	if w.SubStatement.Table != "" {
		subStmt, _ := w.SubStatement.Select.Parse(dialect)
		subWhereStmt := w.SubStatement.WhereStatements.Parse(dialect)
		if w.Op == OperatorExists {
			return fmt.Sprintf("%s (%s%s)", w.Op, subStmt, subWhereStmt)
		}
		return fmt.Sprintf("%s%s%s %s (%s%s)", dialect.GetColumnQuoteLeft(), w.Field, dialect.GetColumnQuoteRight(), w.Op, subStmt, subWhereStmt)
	}

	if strings.Contains(w.Field, ".") {
		return fmt.Sprintf("%s %s %s", pkg.ColumnSplitter(w.Field, dialect.GetColumnQuoteLeft(), dialect.GetColumnQuoteRight()), w.Op, dialect.GetDelimiter())
	}

	return fmt.Sprintf("%s%s%s %s %s", dialect.GetColumnQuoteLeft(), w.Field, dialect.GetColumnQuoteRight(), w.Op, dialect.GetDelimiter())
}

func (ws *WhereStatements) ParseWhereStatements(dialect SQLDialector) string {
	stmt := ""
	for i, v := range ws.Where {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse(dialect)
		if v.Value != nil {
			ws.Values = append(ws.Values, v.Value)
		}
	}
	return stmt
}

func (ws *WhereStatements) ParseWhereInStatements(dialect SQLDialector) string {
	stmt := ""
	if len(ws.Where) > 0 {
		for _, v := range ws.WhereIn {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}
	}

	if len(ws.WhereIn) > 0 {
		for _, v := range ws.WhereIn {
			stmt += v.Parse(dialect)
			ws.Values = append(ws.Values, v.Values...)
		}
	}
	return stmt
}

func (ws *WhereStatements) ParseWhereNotInStatements(dialect SQLDialector) string {
	stmt := ""
	for i, v := range ws.WhereNotIn {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse(dialect)
		ws.Values = append(ws.Values, v.Values...)
	}
	return stmt
}

func (ws *WhereStatements) ParseWhereBetweenStatements(dialect SQLDialector) string {
	stmt := ""
	for i, v := range ws.WhereBetween {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse(dialect)
		ws.Values = append(ws.Values, v.Start, v.End)
	}
	return stmt
}

func (ws *WhereStatements) ParseWhereNotBetweenStatements(dialect SQLDialector) string {
	stmt := ""
	for i, v := range ws.WhereNotBetween {
		if i >= 1 {
			stmt += fmt.Sprintf(" %s ", v.Conj)
		}

		stmt += v.Parse(dialect)
		ws.Values = append(ws.Values, v.Start, v.End)
	}
	return stmt
}

func (w *WhereStatements) Parse(dialect SQLDialector) string {
	stmt := ""
	if len(w.Where) > 0 || len(w.WhereIn) > 0 || len(w.WhereNotIn) > 0 || len(w.WhereBetween) > 0 || len(w.WhereNotBetween) > 0 {
		stmt += " WHERE "
	}

	stmt += w.ParseWhereStatements(dialect)

	stmt += w.ParseWhereInStatements(dialect)

	stmt += w.ParseWhereNotInStatements(dialect)

	stmt += w.ParseWhereBetweenStatements(dialect)

	stmt += w.ParseWhereNotBetweenStatements(dialect)

	return stmt
}
