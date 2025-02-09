package sqlbuilder

import (
	"database/sql"
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
	Where() Builder
	WhereIn() Builder
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

func (b *SQLBuilder) Where() Builder {
	return b
}

func (b *SQLBuilder) WhereIn() Builder {
	return b
}

func NewSQLBuilder(dialect string, sql *sql.DB) *SQLBuilder {
	return &SQLBuilder{dialect: dialect, sql: sql}
}
