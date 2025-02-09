package sqlbuilder

import (
	"database/sql"
	"strings"
)

type SQLBuilder struct {
	dialect string
	sql *sql.DB
	statement string
}

type Builder interface {
	Table(table string) Builder
	Select(cols []string) Builder
	Where() Builder
	WhereIn() Builder
}

func (b *SQLBuilder) Table(table string) Builder {
	b.statement += " FROM " + table
	return b
}

func (b *SQLBuilder) Select(cols []string) Builder {
	b.statement += "SELECT " + strings.Join(cols, ", ")
	return b
}

func (b *SQLBuilder) Where() Builder {
	return b
}

func (b *SQLBuilder) WhereIn() Builder {
	return b
}

func NewSQLBuilder(dialect string, sql *sql.DB) Builder {
	return &SQLBuilder{dialect: dialect, sql: sql}
}
