package sqlbuilder

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

const (
	Eq = "="
	Lte = "<="
	Lt = "<"
	Gte = ">="
	Gt= ">"
)

type SQLBuilder struct {
	dialect string
	sql *sql.DB
	Statement string
}

type Builder interface {
	Table(table string) Builder
	Select(cols []string) Builder
	Where(column string, op string, val interface{}) Builder
	WhereIn(column string, d []string) Builder
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

func (b *SQLBuilder) Where(column string, op string, val interface{}) Builder {
	clause := "WHERE"
	valueBind := "%v"

	if hasWhere(b.Statement) {
		clause = "AND"
	}

	t := reflect.ValueOf(val)
	if t.Kind() == reflect.String {
		valueBind = "'%v'"
	}

	b.Statement += fmt.Sprintf(" %s %s %s " + valueBind, clause, column, op, val)
	return b
}

func (b *SQLBuilder) WhereIn(column string, d []string) Builder {
	clause := "WHERE"
	if hasWhere(b.Statement) {
		clause = "AND"
	}

	b.Statement += fmt.Sprintf(" %s %s in(%s)", clause, column, strings.Join(d, ","))
	return b
}

func NewSQLBuilder(dialect string, sql *sql.DB) *SQLBuilder {
	return &SQLBuilder{dialect: dialect, sql: sql}
}

func hasWhere(statement string) bool {
	return strings.Contains(statement, "WHERE")
}
