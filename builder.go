package sqlbuilder

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

const (
	Eq   = "="
	Lte  = "<="
	Lt   = "<"
	Gte  = ">="
	Gt   = ">"
	Desc = "DESC"
	Asc  = "ASC"
)

const (
	whereOr  = "OR"
	whereAnd = "AND"
)

type SQLBuilder struct {
	dialect   string
	sql       *sql.DB
	hasWhere  bool
	arguments []interface{}
	Statement string
}

type Builder interface {
	Table(table string) Builder
	Select(cols []string) Builder
	Where(column string, comp string, val interface{}) Builder
	WhereIn(column string, d interface{}) Builder
	WhereBetween(column string, start interface{}, end interface{}) Builder
	OrWhere(column string, comp string, val interface{}) Builder
	OrWhereIn(column string, d interface{}) Builder
	OrWhereBetween(column string, start interface{}, end interface{}) Builder
	OrderBy(column string, dir string) Builder
	GroupBy(columns ...string) Builder
	GetSql() string
}

func (b *SQLBuilder) Table(table string) Builder {
	b.Statement += " FROM " + table
	return b
}

func (b *SQLBuilder) Select(cols []string) Builder {
	b.Statement += "SELECT " + strings.Join(cols, ", ")
	return b
}

func (b *SQLBuilder) GetSql() string {
	return b.Statement
}

func (b *SQLBuilder) buildWhere(column string, comp string, val interface{}) Builder {
	b.Statement += fmt.Sprintf(" %s %s ?", column, comp)
	b.arguments = append(b.arguments, val)

	return b
}

func (b *SQLBuilder) buildWhereIn(column string, values interface{}) Builder {
	var inValues string

	ref := reflect.ValueOf(values)
	if ref.Kind() != reflect.Slice {
		return b
	}
	lenValues := ref.Len()

	for i := 0; i < lenValues; i++ {
		val := ref.Index(i).Interface()

		con := "?"

		if i < (lenValues - 1) {
			con += ", "
		}

		inValues += con
		b.arguments = append(b.arguments, val)
	}

	b.Statement += fmt.Sprintf(" %s IN(%s)", column, inValues)
	return b
}

func (b *SQLBuilder) buildWhereBetween(column string, start interface{}, end interface{}) Builder {
	startRef := reflect.ValueOf(start)
	endRef := reflect.ValueOf(end)

	if startRef.Kind() == reflect.String {
		start = fmt.Sprintf("'%s'", startRef.Interface())
	}

	if endRef.Kind() == reflect.String {
		end = fmt.Sprintf("'%s'", endRef.Interface())
	}

	b.Statement += fmt.Sprintf(" %s BETWEEN ? AND ?", column)
	b.arguments = append(b.arguments, start)
	b.arguments = append(b.arguments, end)
	return b
}

func (b *SQLBuilder) Where(column string, comp string, val interface{}) Builder {
	b.setWhereOperator(whereAnd)
	return b.buildWhere(column, comp, val)
}

func (b *SQLBuilder) OrWhere(column string, comp string, val interface{}) Builder {
	b.setWhereOperator(whereOr)
	return b.buildWhere(column, comp, val)
}

func (b *SQLBuilder) WhereIn(column string, values interface{}) Builder {
	b.setWhereOperator(whereAnd)
	return b.buildWhereIn(column, values)
}
func (b *SQLBuilder) OrWhereIn(column string, values interface{}) Builder {
	b.setWhereOperator(whereOr)
	return b.buildWhereIn(column, values)
}

func (b *SQLBuilder) WhereBetween(column string, start interface{}, end interface{}) Builder {
	b.setWhereOperator(whereAnd)
	return b.buildWhereBetween(column, start, end)
}
func (b *SQLBuilder) OrWhereBetween(column string, start interface{}, end interface{}) Builder {
	b.setWhereOperator(whereOr)
	return b.buildWhereBetween(column, start, end)
}
func (b *SQLBuilder) OrderBy(column string, dir string) Builder {
	b.Statement += fmt.Sprintf(" ORDER BY %s %s", column, dir)
	return b
}
func (b *SQLBuilder) GroupBy(column ...string) Builder {
	b.Statement += fmt.Sprintf(" GROUP BY %s", strings.Join(column, ", "))
	return b
}

func (b *SQLBuilder) setWhereOperator(op string) {
	if !b.HasWhere() {
		b.Statement += " WHERE"
	}

	if b.HasWhere() {
		b.Statement += fmt.Sprintf(" %s", op)
	}

	if strings.Contains(b.Statement, "WHERE") {
		b.hasWhere = true
	}
}

func NewSQLBuilder(dialect string, sql *sql.DB) *SQLBuilder {
	return &SQLBuilder{dialect: dialect, sql: sql}
}

func (b *SQLBuilder) HasWhere() bool {
	return b.hasWhere
}
