package sqlbuilder

import "database/sql"

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

func NewSQLBuilder(dialect string, sql *sql.DB) Builder {
	return &SQLBuilder{dialect: dialect, sql: sql}
}
