package sqlbuilder

import (
	"database/sql"
	"fmt"
	"strings"
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
	b.Statement += fmt.Sprintf(" WHERE %s %s '%v'", column, op, val)
	return b
}

func (b *SQLBuilder) WhereIn(column string, d []string) Builder {
	b.Statement += fmt.Sprintf(" WHERE %s in(%s)", column, strings.Join(d, ","))
	return b
}

func NewSQLBuilder(dialect string, sql *sql.DB) *SQLBuilder {
	return &SQLBuilder{dialect: dialect, sql: sql}
}
