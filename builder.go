package sqlbuilder

import (
	"context"
	"database/sql"
	"fmt"
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
	Dialect   string
	sql       *sql.DB
	hasWhere  bool
	arguments []interface{}
	Statement string
}

type Builder interface {
	Table(table string, columns ...string) *SQLBuilder
	Where(statement string, vars ...interface{}) *SQLBuilder
	WhereFunc(statement string, b func(b Builder) *SQLBuilder) *SQLBuilder
	Join(table string, first string, operator string, second string)
	LeftJoin(table string, first string, operator string, second string)
	RightJoin(table string, first string, operator string, second string)
	OrderBy(column string, dir string) *SQLBuilder
	GroupBy(columns ...string) *SQLBuilder
	Limit(n int) *SQLBuilder
	Offset(n int) *SQLBuilder
}

func NewSelect(dialect string, db *sql.DB) *SQLBuilder {
	return &SQLBuilder{
		Dialect:   dialect,
		sql:       db,
		Statement: "SELECT ",
	}
}

func (b *SQLBuilder) Table(table string, columns ...string) *SQLBuilder {
	b.Statement += fmt.Sprintf("%s FROM %s", strings.Join(columns, ","), table)
	return b
}

func (b *SQLBuilder) GetSql() string {
	return b.Statement
}

func (b *SQLBuilder) GetArguments() []interface{} {
	return b.arguments
}

func (b *SQLBuilder) Where(statement string, vars ...interface{}) *SQLBuilder {
	b.Statement += fmt.Sprintf(" WHERE %s", statement)
	b.arguments = append(b.arguments, vars...)
	return b
}

func (b *SQLBuilder) Join(table string, first string, operator string, second string) {
	b.Statement += fmt.Sprintf(" INNER JOIN %s ON %s %s %s", table, first, operator, second)
}

func (b *SQLBuilder) LeftJoin(table string, first string, operator string, second string) {
	b.Statement += fmt.Sprintf(" LEFT JOIN %s ON %s %s %s", table, first, operator, second)
}

func (b *SQLBuilder) RightJoin(table string, first string, operator string, second string) {
	b.Statement += fmt.Sprintf(" RIGHT JOIN %s ON %s %s %s", table, first, operator, second)
}

func (b *SQLBuilder) WhereFunc(statement string, builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(NewSelect(b.Dialect, b.sql))

	b.Statement += fmt.Sprintf(" WHERE %s", statement)
	b.Statement += fmt.Sprintf("(%s)", newBuilder.GetSql())
	b.arguments = append(b.arguments, newBuilder.GetArguments()...)
	return b
}

func (b *SQLBuilder) WhereExists(builder func(b Builder) *SQLBuilder) *SQLBuilder {
	newBuilder := builder(NewSelect(b.Dialect, b.sql))
	b.Statement += fmt.Sprintf(" WHERE EXISTS (%s)", newBuilder.GetSql())
	b.arguments = append(b.arguments, newBuilder.GetArguments()...)

	return b
}

func (b *SQLBuilder) OrderBy(column string, dir string) *SQLBuilder {
	b.Statement += fmt.Sprintf(" ORDER BY %s %s", column, dir)
	return b
}
func (b *SQLBuilder) GroupBy(column ...string) *SQLBuilder {
	b.Statement += fmt.Sprintf(" GROUP BY %s", strings.Join(column, ", "))
	return b
}
func (b *SQLBuilder) Limit(n int) *SQLBuilder {
	b.Statement += fmt.Sprintf(" LIMIT %d", n)
	return b
}
func (b *SQLBuilder) Offset(n int) *SQLBuilder {
	b.Statement += fmt.Sprintf(" OFFSET %d", n)
	return b
}

func (b *SQLBuilder) Find(d interface{}, ctx context.Context) error {
	rows, err := b.runQuery(ctx)
	defer rows.Close()

	if err != nil {
		return err
	}

	if rows.Next() {
		return ScanStruct(d, rows)
	}

	return nil

}
func (b *SQLBuilder) Get(d interface{}, ctx context.Context) error {
	rows, err := b.runQuery(ctx)
	defer rows.Close()

	if err != nil {
		return err
	}

	if err = ScanAll(d, rows); err != nil {
		return err
	}

	return nil
}
func (b *SQLBuilder) Exec(ctx context.Context) error {
	_, err := b.sql.ExecContext(ctx, b.Statement, b.arguments...)

	return err
}

func (b *SQLBuilder) runQuery(ctx context.Context) (*sql.Rows, error) {
	return b.sql.QueryContext(ctx, b.Statement, b.arguments...)
}
